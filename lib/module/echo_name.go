package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

func echoName(cfg *configure.Config, in *message.InboundMsg, actions *actions) bool {
	if !strings.HasPrefix(in.Msg, fmt.Sprintf("%s!", cfg.User.Nick)) {
		return false
	}
	actions.say(fmt.Sprintf("%s!", in.Event.Nick))
	return true
}
