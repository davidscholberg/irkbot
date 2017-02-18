package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

var swears []string

func ConfigInsult(cfg *configure.Config) {
	// initialize swear array
	swearBytes, err := ioutil.ReadFile(cfg.Modules["insult"]["insult_swearfile"])
	if err == nil {
		swears = strings.Split(string(swearBytes), "\n")
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}

func HelpInsult() []string {
	s := "insult [insultee] - insult the given insultee (or self if none" +
		" given)"
	return []string{s}
}

func Insult(p *message.Privmsg) {
	if len(swears) == 0 {
		message.Say(p, "error: no swears")
		return
	}

	insultee := p.Event.Nick
	if len(p.MsgArgs) > 1 {
		insultee = strings.Join(p.MsgArgs[1:], " ")
	}

	response := fmt.Sprintf(
		"%s: you %s %s",
		insultee,
		swears[rand.Intn(len(swears))],
		swears[rand.Intn(len(swears))])

	message.Say(p, response)
}
