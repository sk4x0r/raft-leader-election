package raft

import (
	"log"
	//"os"
	"encoding/json"
	//zmq "github.com/pebbe/zmq4"
	"io/ioutil"
	"strconv"
)
func parseConfigFile(configFile string) Config {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Error parsing the config file")
		panic(err)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Println("Error parsing the config file")
		panic(err)
	}
	return conf
}

//used to create json object of term for storing on disk
type TermJson struct{
	Term int
}


func loadTermFromDisk(serverId int) int{
	var term int
	fileName:=strconv.Itoa(serverId)+".term"
	fileBytes, err := ioutil.ReadFile(fileName)
	if err!=nil{
		//log.Println("Unable to find term stored on disk")
	}
	var termJson TermJson
	err = json.Unmarshal(fileBytes, &termJson)
	if err != nil {
		//log.Print("Error while unmarshalling. Initializing the term to zero")
		term=0
	}else{
		term=termJson.Term
	}
	return term
}

func loadServer(serverId int, conf Config) Server {
	//log.Println("Inside loadServer()")
	var s Server
	s.pid = serverId
	s.peers = conf.getPeers(serverId)
	s.inbox = make(chan *Envelope, 100)
	s.outbox = make(chan *Envelope, 100)
	s.port = conf.getPort(serverId)
	s.peerInfo = conf.getPeerInfo(serverId)
	s.killFollower=make(chan bool, 1)
	s.killCandidate=make(chan bool, 1)
	s.killLeader=make(chan bool, 1)
	s.stopInbox = make(chan bool, 1)
	s.stopOutbox = make(chan bool, 1)
	s.voteRequestReceived  = make(chan bool, 1)
	s.voteReceived  = make(chan bool, 1)
	s.heartbeat  = make(chan int, 1)
	s.alreadyVoted = make(map[int] bool)
	s.steppedDownFromLeadership = make(chan bool, 1)
	
	//raft related
	s.timeout=conf.Timeout
	s.amILeader=false
	s.myTerm=loadTermFromDisk(serverId)
	log.Println(s.Pid(),"Booting up, term=",s.Term())
	go s.handleInbox()
	go s.handleOutbox()
	go s.actAsFollower()
	return s
}
func New(serverId int, configFile string) Server {
	conf := parseConfigFile(configFile)
	s := loadServer(serverId, conf)
	return s
}
