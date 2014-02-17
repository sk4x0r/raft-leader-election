package raft

type Envelope struct {
	Pid   int
	MsgId int64
	Msg   interface{}
}

func (e *Envelope) getPid() int {
	return e.Pid
}

func (e *Envelope) getMsgId() int64 {
	return e.MsgId
}

func (e *Envelope) getMsg() interface{} {
	return e.Msg
}

func (e *Envelope) setPeerId(pid int) {
	e.Pid = pid
}
func (e *Envelope) setMsgId(msgId int64) {
	e.MsgId = msgId
}
func (e *Envelope) setMsg(msg interface{}) {
	e.Msg = msg
}
