package raft

import (
	"log"
    "math/rand"
    "math"
	"time"
	"io/ioutil"
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"strconv"
)

type Peer struct {
	Pid  int
	Ip   string
	Port int
	soc  *zmq.Socket
}

//struct used to create json object of control messages
//which are used for communication among the servers
type ControlMessage struct{
	Command string
	Term int
	ServerId int
}

//struct used to marshal/unmarshal json objects of configuration file
type Config struct {
	Timeout int
	Peers []Peer
}

func (c *Config) getPeers(myId int) []int {
	l := len(c.Peers)
	pids := make([]int, l-1)

	i := 0
	for _, peer := range c.Peers {
		if peer.Pid != myId {
			pids[i] = peer.Pid
			i = i + 1
		}
	}
	return pids
}

func (c *Config) getPort(myId int) int {
	for _, peer := range c.Peers {
		if peer.Pid == myId {
			return peer.Port
		}
	}
	panic("myId" + string(myId) + "doesn't exist")
}

func (c *Config) getPeerInfo(myId int) map[int]Peer {
	peerInfo := make(map[int]Peer)
	for _, peer := range c.Peers {
		if peer.Pid != myId {
			peerInfo[peer.Pid] = peer
		}
	}
	return peerInfo
}

type Server struct {
	pid         int
	peers       []int
	outbox      chan *Envelope
	inbox       chan *Envelope
	stopInbox   chan bool
	stopOutbox  chan bool
	killFollower chan bool
	killCandidate chan bool
	killLeader chan bool
	voteRequestReceived chan bool
	voteReceived chan bool
	steppedDownFromLeadership chan bool
	heartbeat chan int
	port        int
	peerInfo    map[int]Peer
	connections map[int]*zmq.Socket
	
	//raft related variables
	myTerm int
	amILeader bool
	leaderId int
	timeout int
	alreadyVoted map[int]bool
}


//kills all the goroutines and stops the server
func (s *Server) StopServer(){
	s.stopInbox<-true
	s.stopOutbox<-true
	s.killFollower<- true
	s.killCandidate<- true
	s.killLeader<- true
	time.Sleep(5*time.Second)
	
	//TODO: close the following channels
	/*
	close(s.inbox)
	close(s.outbox)
	close(s.killFollower)
	close(s.killCandidate)
	close(s.killLeader)
	close(s.stopInbox)
	close(s.stopOutbox)
	close(s.voteRequestReceived)
	close(s.voteReceived)
	close(s.heartbeat)
	close(s.steppedDownFromLeadership)
	*/
	log.Println(s.Pid(),"Shutting down")
}

func (s *Server) Term() int{
	return s.myTerm
}

func (s *Server) isLeader() bool{
	return s.amILeader
}


func (s *Server) Pid() int {
	return s.pid
}

func (s *Server) Peers() []int {
	return s.peers
}

func (s *Server) Outbox() chan *Envelope {
	return s.outbox
}

func (s *Server) Inbox() chan *Envelope {
	return s.inbox
}

func (s *Server) Port() int {
	return s.port
}


func (s *Server) createVoteRequest() Envelope{
	//create a broadcast message asking for vote
	//log.Println(s.Pid(),"Creating vote request")
	rand.Seed( time.Now().UnixNano())
	msgId:=rand.Int63n(math.MaxInt64)
	msg,err:=json.Marshal(ControlMessage{Command:"REQUEST_VOTE", Term:s.Term(), ServerId:s.Pid()})
	if err !=nil{
		log.Println("Unable to encode Control Message")
		panic(err)
	}
	e := Envelope{Pid: BROADCAST, MsgId: msgId, Msg: string(msg)}
	return e
}


func (s *Server) waitForElectionResult(winsElection chan bool, killWaitingForResult chan bool){
	//log.Println(s.Pid(),"Waiting for election result")
	
	majority:=(len(s.Peers())+1)/2+1
	votecount:=1
	for votecount<majority{
		select{
			case <-s.voteReceived:
			votecount+=1
			//log.Println(s.Pid(),"Vote Received. Votecount=",votecount, "majority=",majority )
				break
			case <-killWaitingForResult:
				return
		}
	}
	//log.Println(s.Pid(),"Got majority")
	winsElection <- true
}


func (s *Server) createHeartbeat() Envelope{
	//create a broadcast message
	
	rand.Seed( time.Now().UnixNano())
	msgId:=rand.Int63n(math.MaxInt64) //TODO: look into uniqueness of msgId
	
	msg,err:=json.Marshal(ControlMessage{Command:"HEARTBEAT", Term:s.Term(), ServerId:s.Pid()})
	if err !=nil{
		log.Fatal("Unable to encode Control Message",err)
	}
	e := Envelope{Pid: BROADCAST, MsgId: msgId, Msg: string(msg)}
	return e
}

func (s *Server)sendHeartbeats(){
	outbox:=s.Outbox()
	heartbeat:=s.createHeartbeat()
	outbox <-&heartbeat
}

//this goroutine keeps executing when server is in leader state
func (s *Server) actAsLeader(){
	s.amILeader=true
	log.Println(s.Pid(),"Leader, term=",s.Term())
	//keep sending heartbeats
	for{
		//log.Println(s.Pid(),"Leader, term=",s.Term())
		select{
			case <-s.killLeader:
				//log.Println(s.Pid(),"killing the leader")
				s.amILeader=false
				log.Println(s.Pid(),"Stepping down from Leadership")
				return
			case leaderTerm:= <-s.heartbeat:
				if leaderTerm>s.myTerm{
					//update my term
					s.updateMyTerm(leaderTerm)
					//step down and become follower again
					//log.Println(s.Pid(),"found somebody with greater term value, stepping down")
					s.amILeader=false
					log.Println(s.Pid(),"Stepping down from Leadership")
					return
				}else if leaderTerm==s.myTerm{
					panic("Two leaders with same term")
				}else{
					//neglect the heartbeat
					break
				}
			//got a voterequest from a higher term, quit
			case <- s.voteRequestReceived:
				s.amILeader=false
				log.Println(s.Pid(),"Stepping down from Leadership")
				s.steppedDownFromLeadership <-true
				return
			case <-time.After(heartbeatInterval()): //TODO: consider storing heartbeat interval in some const/variable
				//log.Println(s.Pid(),"sending heartbeat")
				time.Sleep(1*time.Second)
				s.sendHeartbeats()
			}
		}
		log.Println(s.Pid(),"Stepping down from Leadership")
		s.amILeader=false
}


func (s *Server)actAsCandidate(){
	//become candidate
	for{
		//increment myTerm
		s.updateMyTerm(s.Term()+1)
		//log.Println(s.Pid(),"I am a candidate now and my new term is", s.myTerm)
		outbox:=s.Outbox()
		
		voteRequest:=s.createVoteRequest()
		
		//be ready to receive election results before sending vote request
		winsElection:= make(chan bool, 1)
		killWaitingForResult := make(chan bool, 1)
		go s.waitForElectionResult(winsElection, killWaitingForResult)
		
		outbox <-&voteRequest
		
		WAIT_FOR_ELECTION_RESULT: //label TODO:avoid it if possible
		select{
			case <-s.killCandidate:
				//log.Println(s.Pid(),"killing the candidate")
				killWaitingForResult <- true
				return
			case <-winsElection:
				//log.Println(s.Pid(),"I won the election, now becoming the leader")
				s.actAsLeader()
				//log.Println(s.Pid(),"Good time passes, stepped down from leadership")
				return

			//TODO: look into this case, make sure that the code considers all the conditions
			case leaderTerm:=<-s.heartbeat:
				//log.Println(s.Pid(),"Received heartbeat while I was a candidate")
				if leaderTerm>=s.myTerm{
					//log.Println(s.Pid(),"Found leader with higher term, returning")
					s.updateMyTerm(leaderTerm)
					//step down and become follower again
					return
				}else{ // my term is greater than the term of server from which I received last heartbeat. So I don't care
					//neglect the heartbeat and wait for election results
					//log.Println(s.Pid(),"leader's term is less than mine, continuing candidature")
					goto WAIT_FOR_ELECTION_RESULT
				}
			//received a vote request from higher term, go back to follower state
			case <- s.voteRequestReceived:
				return

			case <-time.After(electionTimeout()):
				//log.Println(s.Pid(),"tried to be a leader, not many voted, going back to be a follower")
				break
			}
		}
}

//This goroutine keeps executing when serever is in a follower state. Same is called during initialization of server
func (s *Server)actAsFollower(){
	for{
		//log.Println(s.Pid(), "- Follower")
		//starts as a follower
		select{
			//if kill signal is received
			case <-s.killFollower:
				//log.Println(s.Pid(),"killing follower")
				return
			//if heartbeat is received
			case <-s.heartbeat:
				//log.Println(s.Pid(),"received heartbeat")
				break
			
			//got a request for vote, consider it as a heartbeat
			case <- s.voteRequestReceived:
				break
			case <-time.After(heartbeatTimeout()):
				//log.Println(s.Pid(),"timed out, going for candidacy")
				s.actAsCandidate()
				break
			}
	}
}


func (s *Server) createVote(serverId int) Envelope{
	rand.Seed( time.Now().UnixNano())
	msgId:=rand.Int63n(math.MaxInt64)
	
	msg,err:=json.Marshal(ControlMessage{Command:"VOTE_RESPONSE"})
	if err !=nil{
		log.Panic("Unable to encode Control Message for voting",err)
	}
	e := Envelope{Pid: serverId, MsgId: msgId, Msg: string(msg)}
	return e
}

func (s *Server) voteForLeader(serverId int){
	//log.Println(s.Pid(),"Sending vote for server", serverId)
	outbox:=s.Outbox()
	msg := s.createVote(serverId)
	outbox <- &msg
}

func(s *Server) updateMyTerm(term int){
	fileName:=strconv.Itoa(s.Pid())+".term"
	
	termJson:=TermJson{term}
	
	fileBytes, err := json.Marshal(termJson)
	if err!=nil{
		panic("Error while marshalling"+err.Error())
	}
	
	err=ioutil.WriteFile(fileName, fileBytes, 0644)
	if err!=nil{
		panic("Error writing to disk"+ err.Error())
	}
	s.myTerm=term
}


func (s *Server)isControlMessage(e Envelope) bool{
	msg, _ :=e.Msg.(string)
	var ctrlMsg ControlMessage
	err := json.Unmarshal([]byte(msg), &ctrlMsg)
	if err != nil {
		return false
	}else{
		//parse control message, and take appropriate action
		if ctrlMsg.Command=="REQUEST_VOTE" {
			if s.myTerm<ctrlMsg.Term{
				//log.Println(s.Pid(),"My term=",s.Term(),"got vote request for term", ctrlMsg.Term)
				s.voteRequestReceived<-true
				if s.isLeader(){//if i am a leader, step down from leadership
					<-s.steppedDownFromLeadership
				}
				s.updateMyTerm(ctrlMsg.Term)
				s.voteForLeader(ctrlMsg.ServerId)
				//log.Println(s.Pid(),"Voted for term",ctrlMsg.Term, "My new term=",s.Term())
				return true
			}else{
				//already voted for this term, discart the request
				return true
			}
		}else if ctrlMsg.Command=="VOTE_RESPONSE" {
			//TODO: check if received vote is for current term or not
			s.voteReceived <- true
			return true
		}else if ctrlMsg.Command=="HEARTBEAT" {
			s.heartbeat <- ctrlMsg.Term
			return true
		}
	}
	return false
}


func (s *Server) handleInbox() {
	//log.Println("handleInbox()")
	responder, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		log.Panicln("Error creating socket",err)
	}
	defer responder.Close()
	bindAddress := "tcp://*:" + strconv.Itoa(s.port)
	responder.Bind(bindAddress)
	responder.SetRcvtimeo(1000*time.Millisecond)
	//log.Println("socket created")
	for {
		select{
			//goroutine should return when something is received on this channel
			case <-s.stopInbox:
				//log.Println(s.Pid(),"stopping inbox")
				return
			//otherwise.. keep receiving and processing requests
			default:
				msg, err := responder.RecvBytes(0)
				//log.Println(s.Pid(), "Received message", msg)
				if err != nil {
					//log.Println("Error receiving message", err.Error())
					break
				}
				envelope := gobToEnvelope(msg)
				if s.isControlMessage(envelope)==false{
					s.inbox <- &envelope
				}
			}
		}
	}

func (s *Server) handleOutbox() {
	s.connections = make(map[int]*zmq.Socket)
	for i := range s.peers {
		peerId := s.peers[i]

		sock, err := zmq.NewSocket(zmq.PUSH)
		if err != nil {
			log.Panicln("Error creating socket", err)
		}
		sockAddr := "tcp://" + s.peerInfo[peerId].Ip + ":" + strconv.Itoa(s.peerInfo[peerId].Port)
		sock.Connect(sockAddr)
		s.connections[peerId] = sock
		//s.peerInfo[peerId].soc= s.connections[peerId]
	}
	for {
		select {
		case message := <-s.outbox:
			envelope := *message
			if envelope.Pid == BROADCAST {
				time.Sleep(50*time.Millisecond)
				for _, conn := range s.connections {
					msg := envelopeToGob(envelope)
					conn.SendBytes(msg, 0)
				}
			} else {
				peerId := envelope.Pid
				conn := s.connections[peerId]
				msg := envelopeToGob(envelope)
				conn.SendBytes(msg, 0)
			}
		case <-s.stopOutbox:
			return
		}
	}
}
