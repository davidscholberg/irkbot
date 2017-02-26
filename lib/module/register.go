package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/thoj/go-ircevent"
	"strings"
)

type CommandModule struct {
	Configure func(*configure.Config)
	GetHelp   func() []string
	Run       func(*message.Privmsg)
}

type ParserModule struct {
	Configure func(*configure.Config)
	GetHelp   func() []string
	Run       func(*message.Privmsg) bool
}

func RegisterModules(conn *irc.Connection, cfg *configure.Config, sayChan chan message.SayMsg) error {
	// register modules
	parserModules := []*ParserModule{}
	cmdMap := make(map[string]*CommandModule)
	for moduleName, _ := range cfg.Modules {
		switch moduleName {
		case "echo_name":
			parserModules = append(parserModules, &ParserModule{ConfigEchoName, nil, EchoName})
		case "help":
			cmdMap["help"] = &CommandModule{ConfigHelp, nil, Help}
		case "slam":
			cmdMap["slam"] = &CommandModule{ConfigSlam, HelpSlam, Slam}
		case "compliment":
			cmdMap["compliment"] = &CommandModule{ConfigCompliment, HelpCompliment, Compliment}
		case "quit":
			cmdMap["quit"] = &CommandModule{ConfigQuit, HelpQuit, Quit}
		case "urban":
			cmdMap["urban"] = &CommandModule{nil, HelpUrban, Urban}
		case "urban_wotd":
			cmdMap["urban_wotd"] = &CommandModule{nil, HelpUrbanWotd, UrbanWotd}
		case "urban_trending":
			cmdMap["urban_trending"] = &CommandModule{nil, HelpUrbanTrending, UrbanTrending}
		case "url":
			parserModules = append(parserModules, &ParserModule{nil, nil, Url})
		default:
			return fmt.Errorf("invalid name '%s' in module config", moduleName)
		}
	}

	for _, m := range cmdMap {
		if m.GetHelp != nil {
			RegisterHelp(m.GetHelp())
		}
		if m.Configure != nil {
			m.Configure(cfg)
		}
	}

	for _, m := range parserModules {
		if m.GetHelp != nil {
			RegisterHelp(m.GetHelp())
		}
		if m.Configure != nil {
			m.Configure(cfg)
		}
	}

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		p := message.Privmsg{}
		p.Msg = e.Message()
		p.MsgArgs = strings.Split(p.Msg, " ")
		p.Dest = e.Arguments[0]
		if !strings.HasPrefix(p.Dest, "#") {
			p.Dest = e.Nick
		}
		p.Event = e
		p.Conn = conn
		p.SayChan = sayChan

		// run parser modules
		for _, m := range parserModules {
			if m.Run(&p) {
				return
			}
		}

		// check commands
		cmdPrefix := cfg.Channel.CmdPrefix
		if cmdPrefix == "" {
			cmdPrefix = "."
		}
		if strings.HasPrefix(p.Msg, cmdPrefix) {
			if m, ok := cmdMap[strings.TrimPrefix(p.MsgArgs[0], cmdPrefix)]; ok {
				m.Run(&p)
			}
		}

	})

	return nil
}
