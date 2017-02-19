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

func ConfigSlam(cfg *configure.Config) {
	// initialize word arrays
	adjectiveBytes, err := ioutil.ReadFile(cfg.Modules["slam"]["adjective_file"])
	if err != nil {
		// TODO: use logger here
		fmt.Fprintln(os.Stderr, err)
		return
	}
	nounBytes, err := ioutil.ReadFile(cfg.Modules["slam"]["noun_file"])
	if err != nil {
		// TODO: use logger here
		fmt.Fprintln(os.Stderr, err)
		return
	}
	adjectives = strings.Split(string(adjectiveBytes), "\n")
	nouns = strings.Split(string(nounBytes), "\n")
}

func HelpSlam() []string {
	s := "slam [victim] - give the victim a verbal smackdown (or self if no victim" +
		" specified)"
	return []string{s}
}

func Slam(p *message.Privmsg) {
	if len(adjectives) == 0 || len(nouns) == 0 {
		message.Say(p, "error: no smackdowns loaded")
		return
	}

	victim := p.Event.Nick
	if len(p.MsgArgs) > 1 {
		victim = strings.TrimSpace(strings.Join(p.MsgArgs[1:], " "))
	}

	response := fmt.Sprintf(
		"%s: u %s %s",
		victim,
		adjectives[rand.Intn(len(adjectives))],
		nouns[rand.Intn(len(nouns))])

	message.Say(p, response)
}
