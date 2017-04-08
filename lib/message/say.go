package message

import (
	"github.com/thoj/go-ircevent"
	"time"
)

type InboundMsg struct {
	Msg     string
	MsgArgs []string
	Src     string
	Event   *irc.Event
}

type OutboundMsg struct {
	Msg  string
	Dest string
	Conn *irc.Connection
}

func SayLoop(outChan chan OutboundMsg) {
	latestMessageTimes := make(map[string]time.Time)

	for o := range outChan {
		timerDuration := time.Duration(0)
		inboundTime := time.Now()

		if lastMessageTime, ok := latestMessageTimes[o.Dest]; ok {
			timerDuration = time.Second - inboundTime.Sub(lastMessageTime)
			if timerDuration < 0 {
				timerDuration = time.Duration(0)
			}
		}

		latestMessageTimes[o.Dest] = inboundTime.Add(timerDuration)
		conn := o.Conn
		dest := o.Dest
		msg := o.Msg
		time.AfterFunc(timerDuration, func() {
			conn.Privmsg(dest, msg)
		})
	}
}
