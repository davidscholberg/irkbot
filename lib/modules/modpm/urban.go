package modpm

import (
    "fmt"
    "strings"
    urbandict "github.com/davidscholberg/go-urbandict"
    "github.com/davidscholberg/irkbot/lib"
)

func Urban(p *lib.Privmsg) bool {
    if ! strings.HasPrefix(p.Msg, "..urban") {
        return false
    }

    var def *urbandict.Definition
    var err error
    nick := p.Event.Nick
    isWotd := strings.HasPrefix(p.Msg, "..urban-wotd")
    if isWotd {
        def, err = urbandict.WordOfTheDay()
        if err != nil {
            p.SayChan <- lib.Say{
                p.Conn,
                p.Dest,
                fmt.Sprintf("%s: %s", nick, err.Error())}
            return true
        }
    } else if len(p.MsgArgs) == 1 {
        def, err = urbandict.Random()
        if err != nil {
            p.SayChan <- lib.Say{
                p.Conn,
                p.Dest,
                fmt.Sprintf("%s: %s", nick, err.Error())}
            return true
        }
    } else {
        def, err = urbandict.Define(strings.Join(p.MsgArgs[1:], " "))
        if err != nil {
            p.SayChan <- lib.Say{
                p.Conn,
                p.Dest,
                fmt.Sprintf("%s: %s", nick, err.Error())}
            return true
        }
    }

    // TODO: implement max message length handling

    if isWotd {
        p.SayChan <- lib.Say{
            p.Conn,
            p.Dest,
            fmt.Sprintf("%s: Word of the day: \"%s\"", nick, def.Word)}
    } else {
        p.SayChan <- lib.Say{
            p.Conn,
            p.Dest,
            fmt.Sprintf("%s: Top definition for \"%s\"", nick, def.Word)}
    }
    for _, line := range strings.Split(def.Definition, "\r\n") {
        p.SayChan <- lib.Say{p.Conn, p.Dest, fmt.Sprintf("%s: %s", nick, line)}
    }
    p.SayChan <- lib.Say{p.Conn, p.Dest, fmt.Sprintf("%s: Example:", nick)}
    for _, line := range strings.Split(def.Example, "\r\n") {
        p.SayChan <- lib.Say{p.Conn, p.Dest, fmt.Sprintf("%s: %s", nick, line)}
    }
    p.SayChan <- lib.Say{
        p.Conn,
        p.Dest,
        fmt.Sprintf("%s: permalink: %s", nick, def.Permalink)}
    return true
}
