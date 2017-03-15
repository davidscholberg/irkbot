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
	messageTimeouts := make(map[string]time.Time)

	for o := range outChan {
		sleepDuration := time.Duration(0)

		if prevTime, ok := messageTimeouts[o.Dest]; ok {
			sleepDuration = time.Second - time.Now().Sub(prevTime)
			if sleepDuration < 0 {
				sleepDuration = time.Duration(0)
			}
		}

		time.Sleep(sleepDuration)
		messageTimeouts[o.Dest] = time.Now()

		o.Conn.Privmsg(o.Dest, o.Msg)
	}
}
