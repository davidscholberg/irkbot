package module

import (
	"fmt"
	"github.com/jholtom/irkbot/lib/configure"
	"github.com/jholtom/irkbot/lib/message"
	"strings"
)

func HelpInterject() []string {
	s := "interject [subject] - give a pedantic rant about the proper way to refer to" +
		" the given subject, defaulting to Linux if no subject given"
	return []string{s}
}

func Interject(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	interMsg := "Linux"
	if len(in.MsgArgs[1:]) > 0 {
		interMsg = strings.Join(in.MsgArgs[1:], " ")
	}
	msg := fmt.Sprintf("I'd just like to interject for a moment. What you're refering"+
		" to as %s, is in fact, GNU *slash* %s, or as I've recently taken to calling it, GNU *plus* %s.", interMsg, interMsg, interMsg)
	actions.Say(msg)
}
