package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
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

	// seed rng
	rand.Seed(time.Now().UnixNano())
}

func HelpCompliment() []string {
	s := "compliment [recipient] - give the recipient a compliment (or self if no recipient" +
		" specified)"
	return []string{s}
}

func Compliment(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if len(compliments) == 0 {
		actions.Say("error: no compliments loaded")
		return
	}

	recipient := in.Event.Nick
	if len(in.MsgArgs) > 1 {
		recipient = strings.TrimSpace(strings.Join(in.MsgArgs[1:], " "))
	}

	response := fmt.Sprintf(
		"%s: %s",
		recipient,
		compliments[rand.Intn(len(compliments))],
	)

	actions.Say(response)
}
