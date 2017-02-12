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
	swearBytes, err := ioutil.ReadFile(cfg.Module.InsultSwearfile)
	if err == nil {
		swears = strings.Split(string(swearBytes), "\n")
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}

func HelpInsult() []string {
	s := "..insult [insultee] - insult the given insultee (or self if none" +
		" given)"
	return []string{s}
}

func Insult(p *message.Privmsg) bool {
	if !strings.HasPrefix(p.Msg, "..insult") {
		return false
	}

	if len(swears) == 0 {
		message.Say(p, "error: no swears")
		return true
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
	return true
}
