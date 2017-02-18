package module

import (
	"github.com/davidscholberg/irkbot/lib/message"
)

func Quit(p *message.Privmsg) {
	p.Conn.Quit()
}
