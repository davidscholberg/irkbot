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

var compliments []string

func ConfigCompliment(cfg *configure.Config) {
	// initialize compliment array
	complimentBytes, err := ioutil.ReadFile(cfg.Modules["compliment"]["file"])
	if err != nil {
		// TODO: use logger here
		fmt.Fprintln(os.Stderr, err)
		return
	}
	compliments = strings.Split(string(complimentBytes), "\n")
}

func HelpCompliment() []string {
	s := "compliment [recipient] - give the recipient a compliment (or self if no recipient" +
		" specified)"
	return []string{s}
}

func Compliment(p *message.Privmsg) {
	if len(compliments) == 0 {
		message.Say(p, "error: no compliments loaded")
		return
	}

	recipient := p.Event.Nick
	if len(p.MsgArgs) > 1 {
		recipient = strings.TrimSpace(strings.Join(p.MsgArgs[1:], " "))
	}

	response := fmt.Sprintf(
		"%s: %s",
		recipient,
		compliments[rand.Intn(len(compliments))],
	)

	message.Say(p, response)
}
