package modpm

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

var swears []string

func ConfigInsult(cfg *lib.Config) {
	// initialize swear array
	swearBytes, err := ioutil.ReadFile(cfg.Module.Insult_swearfile)
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

func Insult(p *lib.Privmsg) bool {
	if !strings.HasPrefix(p.Msg, "..insult") {
		return false
	}

	if len(swears) == 0 {
		lib.Say(p, "error: no swears")
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

	lib.Say(p, response)
	return true
}
