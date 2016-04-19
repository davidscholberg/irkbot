package main

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib"
	"github.com/davidscholberg/irkbot/lib/modules/modpm"
	goirc "github.com/thoj/go-ircevent"
	gcfg "gopkg.in/gcfg.v1"
	"os"
	"strings"
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

	conn.VerboseCallbackHandler = cfg.Connection.Verbose_callback_handler
	conn.Debug = cfg.Connection.Debug

	conn.AddCallback("001", func(e *goirc.Event) {
		conn.Join(cfg.Channel.Channelname)
	})

	conn.AddCallback("366", func(e *goirc.Event) {
		if len(cfg.Channel.Greeting) != 0 {
			conn.Privmsg(e.Arguments[1], cfg.Channel.Greeting)
		}
	})

	// register modules
	var pmMods []*lib.Module
	modpm.RegisterMods(func(m *lib.Module) {
		pmMods = append(pmMods, m)
		if m.Configure != nil {
			m.Configure(&cfg)
		}
	})

	// TODO: start multiple sayLoops, one per conn
	// TODO: pass conn to sayLoop instead of privmsg callbacks?
	sayChan := make(chan lib.SayMsg)
	go lib.SayLoop(sayChan)

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

		for _, mod := range pmMods {
			if mod.Run(&p) {
				break
			}
		}
	})

	// TODO: add time-based modules

	conn.Loop()
}
