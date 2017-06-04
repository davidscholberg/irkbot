package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dvdmuckle/irkbot/lib/configure"
	"github.com/dvdmuckle/irkbot/lib/message"
	"net/http"
	"os"
	"strings"
)

type doomStruct struct {
	Type string `json:"type"`
}

var doomHost string
var doomValids = []string{"shoot", "forward", "backward", "left", "right", "use"}

func ConfigDoom(cfg *configure.Config) {
	doomHost = cfg.Modules["doom"]["doom_host"]
}

func HelpDoom() []string {
	s := "doom <command> - play doom!"
	return []string{s}
}

//Sanitize input
func Doom(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if len(in.MsgArgs[1:]) == 0 {
		actions.Say("enter a command, dipstick")
		return
	}
	doomCommand := strings.Join(in.MsgArgs[1:], " ")
	doomValid := false
	for _, v := range doomValids {
		if doomCommand == v {
			doomValid = true
			break
		}
	}
	if !doomValid {
		actions.Say(fmt.Sprintf("invalid command, commands are "+"%s", strings.Join(doomValids, ", ")))
		return
	}
	doomToPost := doomStruct{Type: doomCommand}
	jsonValue, err := json.Marshal(doomToPost)
	if err != nil {
		// handle err
		fmt.Fprintln(os.Stderr, err)
		actions.Say("something borked, try again")
		return
	}
	resp, err := http.Post(doomHost, "application/json", bytes.NewReader(jsonValue))
	if err != nil {
		// handle err
		fmt.Fprintln(os.Stderr, err)
		actions.Say("something borked, try again")
		return
	}
	defer resp.Body.Close()
}
