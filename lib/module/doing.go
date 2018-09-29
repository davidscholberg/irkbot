package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"math/rand"
	"strings"
	"time"
)

func ConfigDoing(cfg *configure.Config) {
	rand.Seed(time.Now().Unix())
}
func HelpDoing() []string {
	s := "doing [subject] - sun is not doing, [subject] is doing; defaults to command invoker"
	return []string{s}
}

func Doing(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	doingNot := [...]string{
		"sun is",
		"stars are",
		"trees are",
		"moon is",
		"planets are",
		"galaxies are"}
	doingSub := in.Event.Nick
	if len(in.MsgArgs[1:]) > 0 {
		doingSub = strings.Join(in.MsgArgs[1:], " ")
	}
	msg := fmt.Sprintf("%s not doing, %s is doing", doingNot[rand.Intn(len(doingNot))], doingSub)
	actions.Say(msg)
}
