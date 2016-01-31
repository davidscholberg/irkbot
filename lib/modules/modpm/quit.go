package modpm

import (
	"github.com/davidscholberg/irkbot/lib"
)

func Quit(p *lib.Privmsg) bool {
	if p.Msg != "..quit" {
		return false
	}
	p.Conn.Quit()
	return true
}
