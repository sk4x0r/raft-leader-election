package raft

import (
	//"flag"
	//"fmt"
	//"testing"
	"math/rand"
	"time"
	"bytes"
	"encoding/gob"
	"log"
	//"github.com/sk4x0r/cluster"
)

const (
	BROADCAST = -1
	PATH_TO_CONFIG = "config.json"
	ELECTION_TIMEOUT = 5000  //in millisecond
	HEARTBEAT_TIMEOUT = 3000 //in millisecond
	HEARTBEAT_INTERVAL = 1000 //in millisecond
)

func heartbeatInterval() time.Duration{
	return time.Duration(HEARTBEAT_INTERVAL)*time.Millisecond
}


//cite:http://blog.golang.org/gobs-of-data
//cite: Pushkar Khadilkar
func gobToEnvelope(gobbed []byte) Envelope {
	buf := bytes.NewBuffer(gobbed)
	dec := gob.NewDecoder(buf)
	var envelope Envelope
	dec.Decode(&envelope)
	return envelope
}

//cite: http://blog.golang.org/gobs-of-data
//cite: Pushkar Khadilkar
func envelopeToGob(envelope Envelope) []byte {
	var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err:=enc.Encode(envelope)
    if err != nil {
        log.Fatal("encode error:", err)
    }
	return buf.Bytes()
}


func heartbeatTimeout() time.Duration{rand.Seed( time.Now().UnixNano())
	return time.Duration((HEARTBEAT_TIMEOUT+rand.Intn(HEARTBEAT_TIMEOUT)))*time.Millisecond
}

func electionTimeout() time.Duration{
	rand.Seed( time.Now().UnixNano())
	return time.Duration((ELECTION_TIMEOUT+rand.Intn(ELECTION_TIMEOUT)))*time.Millisecond
}

func createDummyMessage(msgId int, pid int) Envelope {
	e := Envelope{Pid: pid, MsgId: int64(msgId), Msg: "Dummy Message"}
	return e
}

func sendMessages(s Server, count int, pid int) {
	outbox := s.Outbox()
	for i := 0; i < count; i++ {
		msg := createDummyMessage(i, pid)
		outbox <- &msg
		time.Sleep(50 * time.Millisecond)
	}
	s.StopServer()
	//fmt.Println("Sent ", count, " messages to ", pid)
}

func receiveMessages(s Server, count int, success chan bool) {
	inbox := s.Inbox()
	for i := 0; i < count; i++ {
		<-inbox
	}
	s.StopServer()
	//fmt.Println("Received ", count, "messages")
	success <- true
}
