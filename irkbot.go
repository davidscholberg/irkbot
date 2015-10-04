package main

import (
    "fmt"
    "io/ioutil"
    "math/rand"
    "os"
    "strings"
    "time"
    goirc "github.com/thoj/go-ircevent"
    urbandict "github.com/davidscholberg/go-urbandict"
    gcfg "gopkg.in/gcfg.v1"
)

type config struct {
    User struct {
        Nick string
        User string
    }
    Server struct {
        Host string
        Port string
    }
    Channel struct {
        Channelname string
        Greeting string
    }
    Module struct {
        Insult_swearfile string
    }
}

type privmsg struct {
    msg string
    msgArgs []string
    dest string
    event *goirc.Event
    conn *goirc.Connection
    sayChan chan say
}

type say struct {
    conn *goirc.Connection
    dest string
    msg string
}

var swears []string

func main() {
    // get config
    confPath := fmt.Sprintf("%s/.config/irkbot/irkbot.ini", os.Getenv("HOME"))
    cfg := config{}
    err := gcfg.ReadFileInto(&cfg, confPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    // initialize swear array
    swearBytes, err := ioutil.ReadFile(cfg.Module.Insult_swearfile)
    if err == nil {
        swears = strings.Split(string(swearBytes), "\n")
    } else {
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

    privmsgCallbacks := []func(*privmsg) bool{
        privmsgEchoName,
        privmsgQuit,
        privmsgInsult,
        privmsgUrban}

    // TODO: start multiple sayLoops, one per conn
    // TODO: pass conn to sayLoop instead of privmsg callbacks?
    sayChan := make(chan say)
    go sayLoop(sayChan)

    conn.AddCallback("PRIVMSG", func(e *goirc.Event) {
        p := privmsg{}
        p.msg = e.Message()
        p.msgArgs = strings.Split(p.msg, " ")
        p.dest = e.Arguments[0]
        if !strings.HasPrefix(p.dest, "#") {
            p.dest = e.Nick
        }
        p.event = e
        p.conn = conn
        p.sayChan = sayChan

        for _, callback := range privmsgCallbacks {
            if callback(&p) {
                break
            }
        }
    })

    conn.Loop()
}

func privmsgEchoName(p *privmsg) bool {
    if p.msg != "irkbot!" {
        return false
    }
    p.sayChan <- say{p.conn, p.dest, fmt.Sprintf("%s!", p.event.Nick)}
    return true
}

func privmsgQuit(p *privmsg) bool {
    if p.msg != "..quit" {
        return false
    }
    p.conn.Quit()
    return true
}

func privmsgInsult(p *privmsg) bool {
    if ! strings.HasPrefix(p.msg, "..insult") {
        return false
    }

    if len(swears) == 0 {
        p.sayChan <- say{p.conn, p.dest, "error: no swears"}
        return true
    }

    insultee := p.event.Nick
    if len(p.msgArgs) > 1 {
        insultee = strings.Join(p.msgArgs[1:], " ")
    }

    response := fmt.Sprintf(
        "%s: you %s %s",
        insultee,
        swears[rand.Intn(len(swears))],
        swears[rand.Intn(len(swears))])

    p.sayChan <- say{p.conn, p.dest, response}
    return true
}

func privmsgUrban(p *privmsg) bool {
    if ! strings.HasPrefix(p.msg, "..urban") {
        return false
    }

    var def *urbandict.Definition
    var err error
    nick := p.event.Nick
    if len(p.msgArgs) == 1 {
        def, err = urbandict.Random()
        if err != nil {
            p.sayChan <- say{
                p.conn,
                p.dest,
                fmt.Sprintf("%s: %s", nick, err.Error())}
            return true
        }
    } else {
        def, err = urbandict.Define(strings.Join(p.msgArgs[1:], " "))
        if err != nil {
            p.sayChan <- say{
                p.conn,
                p.dest,
                fmt.Sprintf("%s: %s", nick, err.Error())}
            return true
        }
    }

    // TODO: implement max message length handling

    p.sayChan <- say{
        p.conn,
        p.dest,
        fmt.Sprintf("%s: Top definition for \"%s\"", nick, def.Word)}
    for _, line := range strings.Split(def.Definition, "\r\n") {
        p.sayChan <- say{p.conn, p.dest, fmt.Sprintf("%s: %s", nick, line)}
    }
    p.sayChan <- say{p.conn, p.dest, fmt.Sprintf("%s: Example:", nick)}
    for _, line := range strings.Split(def.Example, "\r\n") {
        p.sayChan <- say{p.conn, p.dest, fmt.Sprintf("%s: %s", nick, line)}
    }
    p.sayChan <- say{
        p.conn,
        p.dest,
        fmt.Sprintf("%s: permalink: %s", nick, def.Permalink)}
    return true
}

func sayLoop(sayChan chan say) {
    sayTimeouts := make(map[string]time.Time)

    for s := range sayChan {
        sleepDuration := time.Duration(0)

        if prevTime, ok := sayTimeouts[s.dest]; ok {
            sleepDuration = time.Second - time.Now().Sub(prevTime)
            if sleepDuration < 0 {
                sleepDuration = time.Duration(0)
            }
        }

        time.Sleep(sleepDuration)
        sayTimeouts[s.dest] = time.Now()

        s.conn.Privmsg(s.dest, s.msg)
    }
}
