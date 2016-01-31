package modpm

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib"
	"strings"
)

var nick string

func ConfigEchoName(cfg *lib.Config) {
	nick = cfg.User.Nick
}

func EchoName(p *lib.Privmsg) bool {
	if !strings.HasPrefix(p.Msg, fmt.Sprintf("%s!", nick)) {
		return false
	}
	lib.Say(p, fmt.Sprintf("%s!", p.Event.Nick))
	return true
}
