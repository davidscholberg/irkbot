package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
)

var owner string
var denyMessage string

func ConfigQuit(cfg *configure.Config) {
	owner = cfg.Admin.Owner
	denyMessage = cfg.Admin.DenyMessage
}

func HelpQuit() []string {
	s := "quit - ragequit IRC (requires owner privilege)"
	return []string{s}
}

func Quit(p *message.Privmsg) {
	if p.Event.Nick != owner {
		message.Say(p, fmt.Sprintf("%s: %s", p.Event.Nick, denyMessage))
		return
	}
	p.Conn.Quit()
}
