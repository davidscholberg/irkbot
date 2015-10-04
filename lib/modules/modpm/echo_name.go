package modpm

import (
    "fmt"
    "github.com/davidscholberg/irkbot/lib"
)

func EchoName(p *lib.Privmsg) bool {
    if p.Msg != "irkbot!" {
        return false
    }
    p.SayChan <- lib.Say{p.Conn, p.Dest, fmt.Sprintf("%s!", p.Event.Nick)}
    return true
}
