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

var adjectives []string
var nouns []string

func ConfigInsult(cfg *configure.Config) {
	// initialize word arrays
	adjectiveBytes, err := ioutil.ReadFile(cfg.Modules["insult"]["adjective_file"])
	if err != nil {
		// TODO: use logger here
		fmt.Fprintln(os.Stderr, err)
		return
	}
	nounBytes, err := ioutil.ReadFile(cfg.Modules["insult"]["noun_file"])
	if err != nil {
		// TODO: use logger here
		fmt.Fprintln(os.Stderr, err)
		return
	}
	adjectives = strings.Split(string(adjectiveBytes), "\n")
	nouns = strings.Split(string(nounBytes), "\n")
}

func HelpInsult() []string {
	s := "insult [insultee] - insult the given insultee (or self if none" +
		" given)"
	return []string{s}
}

func Insult(p *message.Privmsg) {
	if len(adjectives) == 0 || len(nouns) == 0 {
		message.Say(p, "error: no insults loaded")
		return
	}

	insultee := p.Event.Nick
	if len(p.MsgArgs) > 1 {
		insultee = strings.Join(p.MsgArgs[1:], " ")
	}

	response := fmt.Sprintf(
		"%s: u %s %s",
		insultee,
		adjectives[rand.Intn(len(adjectives))],
		nouns[rand.Intn(len(nouns))])

	message.Say(p, response)
}
