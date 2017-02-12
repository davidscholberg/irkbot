package module

import (
	"github.com/davidscholberg/irkbot/lib/message"
)

func Quit(p *message.Privmsg) bool {
	if p.Msg != "..quit" {
		return false
	}
	p.Conn.Quit()
	return true
}
