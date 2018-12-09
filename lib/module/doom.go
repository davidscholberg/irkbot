package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"os"
	"strings"
)

type doomStruct struct {
	Type string `json:"type"`
}

var doomValids = []string{"shoot", "forward", "backward", "left", "right", "use"}

func helpDoom() []string {
	s := "doom <command> - play doom!"
	return []string{s}
}

func doom(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if len(in.MsgArgs[1:]) == 0 {
		actions.say("enter a command plz")
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
		actions.say(fmt.Sprintf("invalid command, commands are "+"%s", strings.Join(doomValids, ", ")))
		return
	}
	doomToPost := doomStruct{Type: doomCommand}
	jsonValue, err := json.Marshal(doomToPost)
	if err != nil {
		// handle err
		fmt.Fprintln(os.Stderr, err)
		actions.say("something borked, try again")
		return
	}
	doomHost := cfg.Modules["doom"]["doom_host"]
	response, err := actions.httpPost(doomHost, "application/json", bytes.NewReader(jsonValue))
	if err != nil {
		// handle err
		fmt.Fprintln(os.Stderr, err)
		actions.say("something borked, try again")
		return
	}
	defer response.Body.Close()
}
