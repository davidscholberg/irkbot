package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

// directMessageLog logs direct messages sent to the bot
func directMessageLog(cfg *configure.Config, in *message.InboundMsg, actions *actions) bool {
	// only handle DMs
	if strings.HasPrefix(in.Src, "#") {
		return false
	}

	fmt.Printf("direct message: <%s> %s\n", in.Event.Nick, in.Msg)

	// don't consume the message, in case there are commands in it
	return false
}
