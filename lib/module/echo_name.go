package module

import (
	"fmt"
	"github.com/jholtom/irkbot/lib/configure"
	"github.com/jholtom/irkbot/lib/message"
	"strings"
)

func EchoName(cfg *configure.Config, in *message.InboundMsg, actions *Actions) bool {
	if !strings.HasPrefix(in.Msg, fmt.Sprintf("%s!", cfg.User.Nick)) {
		return false
	}
	actions.Say(fmt.Sprintf("%s!", in.Event.Nick))
	return true
}
