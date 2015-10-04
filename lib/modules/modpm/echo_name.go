package modpm

import (
    "fmt"
    "strings"
    "github.com/davidscholberg/irkbot/lib"
)

var nick string

func ConfigEchoName(cfg *lib.Config) {
    nick = cfg.User.Nick
}

func EchoName(p *lib.Privmsg) bool {
    if ! strings.HasPrefix(p.Msg, fmt.Sprintf("%s!", nick)) {
        return false
    }
    p.SayChan <- lib.Say{p.Conn, p.Dest, fmt.Sprintf("%s!", p.Event.Nick)}
    return true
}
