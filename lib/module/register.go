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
	Run       func(*message.InboundMsg, *Actions)
}

type ParserModule struct {
	Configure func(*configure.Config)
	GetHelp   func() []string
	Run       func(*message.InboundMsg, *Actions) bool
}

type Actions struct {
	Quit func()
	Say  func(string)
}

func RegisterModules(conn *irc.Connection, cfg *configure.Config, outChan chan message.OutboundMsg) error {
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
		case "quote":
			parserModules = append(parserModules, &ParserModule{ConfigQuote, nil, UpdateQuoteBuffer})
			cmdMap["grab"] = &CommandModule{nil, HelpGrabQuote, GrabQuote}
			cmdMap["quote"] = &CommandModule{nil, HelpGetQuote, GetQuote}
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
		inboundMsg := message.InboundMsg{}
		inboundMsg.Msg = e.Message()
		inboundMsg.MsgArgs = strings.Fields(inboundMsg.Msg)
		inboundMsg.Src = e.Arguments[0]
		if !strings.HasPrefix(inboundMsg.Src, "#") {
			inboundMsg.Src = e.Nick
		}
		inboundMsg.Event = e

		outboundMsg := message.OutboundMsg{}
		outboundMsg.Dest = inboundMsg.Src
		outboundMsg.Conn = conn
		//p.SayChan = sayChan

		sayFunc := func(msg string) {
			outboundMsg.Msg = msg
			outChan <- outboundMsg
		}
		quitFunc := func() {
			conn.Quit()
		}
		actions := Actions{
			Quit: quitFunc,
			Say:  sayFunc,
		}

		// run parser modules
		for _, m := range parserModules {
			if m.Run(&inboundMsg, &actions) {
				return
			}
		}

		// check commands
		cmdPrefix := cfg.Channel.CmdPrefix
		if cmdPrefix == "" {
			cmdPrefix = "."
		}
		if strings.HasPrefix(inboundMsg.Msg, cmdPrefix) {
			if m, ok := cmdMap[strings.TrimPrefix(inboundMsg.MsgArgs[0], cmdPrefix)]; ok {
				m.Run(&inboundMsg, &actions)
			}
		}

	})

	return nil
}
