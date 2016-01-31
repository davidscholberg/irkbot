package lib

import (
	"time"
)

func Say(p *Privmsg, msg string) {
	p.SayChan <- SayMsg{p.Conn, p.Dest, msg}
}

func SayLoop(sayChan chan SayMsg) {
	sayTimeouts := make(map[string]time.Time)

	for s := range sayChan {
		sleepDuration := time.Duration(0)

		if prevTime, ok := sayTimeouts[s.Dest]; ok {
			sleepDuration = time.Second - time.Now().Sub(prevTime)
			if sleepDuration < 0 {
				sleepDuration = time.Duration(0)
			}
		}

		time.Sleep(sleepDuration)
		sayTimeouts[s.Dest] = time.Now()

		s.Conn.Privmsg(s.Dest, s.Msg)
	}
}
