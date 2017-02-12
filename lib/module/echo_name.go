package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

var nick string

func ConfigEchoName(cfg *configure.Config) {
	nick = cfg.User.Nick
}

func EchoName(p *message.Privmsg) bool {
	if !strings.HasPrefix(p.Msg, fmt.Sprintf("%s!", nick)) {
		return false
	}
	message.Say(p, fmt.Sprintf("%s!", p.Event.Nick))
	return true
}
