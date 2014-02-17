package raft

import (
	//"flag"
	//"fmt"
	"testing"
	"time"
	"fmt"
	//"github.com/sk4x0r/cluster"
)

func TestFiveServers(t *testing.T) {
	fmt.Println("Testing for five servers")
	
	//start three servers
	s1:=New(1001, PATH_TO_CONFIG)
	s2:=New(1002, PATH_TO_CONFIG)
	s3:=New(1003, PATH_TO_CONFIG)
	s4:=New(1004, PATH_TO_CONFIG)
	s5:=New(1005, PATH_TO_CONFIG)

	time.Sleep(5*time.Second)
	
	//wait for a while
	
	//s1 crashes
	s1.StopServer()
	s1=New(1001, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	
	//s2 crashes
	s2.StopServer()
	s2=New(1002, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//s3 crashes
	s3.StopServer()
	s3=New(1003, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//s4 crashes
	s4.StopServer()
	s4=New(1004, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//s5 crashes
	s5.StopServer()
	s5=New(1005, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	
	//s1,s2 crash
	s1.StopServer()
	s2.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s1=New(1001, PATH_TO_CONFIG)
	s2=New(1002, PATH_TO_CONFIG)

	time.Sleep(5*time.Second)
	
	
	//s2,s3 crash
	s2.StopServer()
	s3.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s2=New(1002, PATH_TO_CONFIG)
	s3=New(1003, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//s4,s5 crash
	s4.StopServer()
	s5.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s4=New(1004, PATH_TO_CONFIG)
	s5=New(1005, PATH_TO_CONFIG)
	
	time.Sleep(5*time.Second)
	
	
	//s1,s3,s4,s5 crash
	s1.StopServer()
	s3.StopServer()
	s4.StopServer()
	s5.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s1=New(1001, PATH_TO_CONFIG)
	s3=New(1003, PATH_TO_CONFIG)
	s4=New(1004, PATH_TO_CONFIG)
	s5=New(1005, PATH_TO_CONFIG)
	
	time.Sleep(5*time.Second)

	s1.StopServer()
	s2.StopServer()
	s3.StopServer()
	s4.StopServer()
	s5.StopServer()	
}

/*
func TestThreeServers(t *testing.T) {
	fmt.Println("Testing for three servers")
	
	//start three servers
	s1:=New(1001, PATH_TO_CONFIG)
	s2:=New(1002, PATH_TO_CONFIG)
	s3:=New(1003, PATH_TO_CONFIG)
		
	time.Sleep(5*time.Second)
	
	//wait for a while
	
	//s1 crashes
	s1.StopServer()
	New(1001, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	
	//s2 crashes
	s2.StopServer()
	New(1002, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//s3 crashes
	s3.StopServer()
	New(1003, PATH_TO_CONFIG)
	
	time.Sleep(20*time.Second)
	
	
	
	//s1,s2 crash
	s1.StopServer()
	s2.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s1=New(1001, PATH_TO_CONFIG)
	s2=New(1002, PATH_TO_CONFIG)

	time.Sleep(5*time.Second)
	
	
	//s2,s3 crash
	s2.StopServer()
	s3.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s2=New(1002, PATH_TO_CONFIG)
	s3=New(1003, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//s1,s3 crash
	s1.StopServer()
	s3.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	s1=New(1001, PATH_TO_CONFIG)
	s3=New(1003, PATH_TO_CONFIG)
	
	time.Sleep(5*time.Second)
	s1.StopServer()
	s2.StopServer()
	s3.StopServer()
	
}



func TestTwoServers(t *testing.T) {
	fmt.Println("Testing two servers")
	
	//start two servers
	s1:=New(1001, PATH_TO_CONFIG)
	s2:=New(1002, PATH_TO_CONFIG)
		
	//wait for a while
	time.Sleep(5*time.Second)
	
	//stop s1
	s1.StopServer()
	
	//wait
	time.Sleep(5*time.Second)
	
	//restart s1
	New(1001, PATH_TO_CONFIG)
	time.Sleep(5*time.Second)
	
	//stop s2
	s2.StopServer()
	time.Sleep(5*time.Second)
	
	//start s2
	New(1002, PATH_TO_CONFIG)
	time.Sleep(5*time.Minute)
	
	s1.StopServer()
	s2.StopServer()
}
*/
