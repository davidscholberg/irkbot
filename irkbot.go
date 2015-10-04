package main

import (
    "fmt"
    "os"
    "strings"
    "time"
    goirc "github.com/thoj/go-ircevent"
    gcfg "gopkg.in/gcfg.v1"
    "github.com/davidscholberg/irkbot/lib"
    "github.com/davidscholberg/irkbot/lib/modules/modpm"
)

func main() {
    // get config
    confPath := fmt.Sprintf("%s/.config/irkbot/irkbot.ini", os.Getenv("HOME"))
    cfg := lib.Config{}
    err := gcfg.ReadFileInto(&cfg, confPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    conn := goirc.IRC(cfg.User.Nick, cfg.User.User)
    err = conn.Connect(fmt.Sprintf(
        "%s:%s",
        cfg.Server.Host,
        cfg.Server.Port))
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    conn.VerboseCallbackHandler = true
    conn.Debug = true

    conn.AddCallback("001", func(e *goirc.Event) {
        conn.Join(cfg.Channel.Channelname)
    })

    conn.AddCallback("366", func(e *goirc.Event) {
        if len(cfg.Channel.Greeting) != 0 {
            conn.Privmsg(e.Arguments[1], cfg.Channel.Greeting)
        }
    })

    // function calls for module-specific configuration
    modpm.ConfigInsult(&cfg)

    privmsgCallbacks := []func(*lib.Privmsg) bool{
        modpm.EchoName,
        modpm.Quit,
        modpm.Insult,
        modpm.Urban}

    // TODO: start multiple sayLoops, one per conn
    // TODO: pass conn to sayLoop instead of privmsg callbacks?
    sayChan := make(chan lib.Say)
    go sayLoop(sayChan)

    conn.AddCallback("PRIVMSG", func(e *goirc.Event) {
        p := lib.Privmsg{}
        p.Msg = e.Message()
        p.MsgArgs = strings.Split(p.Msg, " ")
        p.Dest = e.Arguments[0]
        if !strings.HasPrefix(p.Dest, "#") {
            p.Dest = e.Nick
        }
        p.Event = e
        p.Conn = conn
        p.SayChan = sayChan

        for _, callback := range privmsgCallbacks {
            if callback(&p) {
                break
            }
        }
    })

    // TODO: add time-based modules

    conn.Loop()
}

func sayLoop(sayChan chan lib.Say) {
    sayTimeouts := make(map[string]time.Time)

    for s := range sayChan {
        sleepDuration := time.Duration(0)

        if prevTime, ok := sayTimeouts[s.Dest]; ok {
            sleepDuration = time.Second - time.Now().Sub(prevTime)
            if sleepDuration < 0 {
                sleepDuration = time.Duration(0)
            }
        }

        time.Sleep(sleepDuration)
        sayTimeouts[s.Dest] = time.Now()

        s.Conn.Privmsg(s.Dest, s.Msg)
    }
}
