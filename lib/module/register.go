package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
)

type CommandModule struct {
	Configure func(*configure.Config)
	GetHelp   func() []string
	Run       func(*configure.Config, *message.InboundMsg, *Actions)
}

type ParserModule struct {
	Configure func(*configure.Config)
	Run       func(*configure.Config, *message.InboundMsg, *Actions) bool
}

type TickerModule struct {
	Configure   func(*configure.Config)
	GetDuration func(*configure.Config) time.Duration
	Run         func(*configure.Config, time.Time, *Actions)
	Ticker      *time.Ticker
}

type Actions struct {
	Quit  func()
	Say   func(string)
	SayTo func(string, string)
}

func RegisterModules(conn *irc.Connection, cfg *configure.Config, outChan chan message.OutboundMsg) error {
	cmdMap := make(map[string]*CommandModule)
	parserModules := []*ParserModule{}
	tickerModules := []*TickerModule{}
	for moduleName, _ := range cfg.Modules {
		switch moduleName {
		case "echo_name":
			parserModules = append(parserModules, &ParserModule{nil, EchoName})
		case "help":
			cmdMap["help"] = &CommandModule{nil, nil, Help}
			parserModules = append(parserModules, &ParserModule{nil, ParseHelp})
		case "slam":
			cmdMap["slam"] = &CommandModule{ConfigSlam, HelpSlam, Slam}
		case "compliment":
			cmdMap["compliment"] = &CommandModule{ConfigCompliment, HelpCompliment, GiveCompliment}
                case "paste":
                        cmdMap["paste"] = &CommandModule{nil, HelpGetPaste, GetPaste}
                        cmdMap["store"] = &CommandModule{nil, HelpStorePaste, StorePaste}
		case "quit":
			cmdMap["quit"] = &CommandModule{nil, HelpQuit, Quit}
		case "quote":
			cmdMap["grab"] = &CommandModule{nil, HelpGrabQuote, GrabQuote}
			cmdMap["quote"] = &CommandModule{nil, HelpGetQuote, GetQuote}
			parserModules = append(parserModules, &ParserModule{ConfigQuote, UpdateQuoteBuffer})
			tickerModules = append(
				tickerModules,
				&TickerModule{
					nil,
					GetCleanQuoteBufferDuration,
					CleanQuoteBuffer,
					nil,
				},
			)
		case "say":
			cmdMap["say"] = &CommandModule{nil, HelpSay, Say}
		case "urban":
			cmdMap["urban"] = &CommandModule{nil, HelpUrban, Urban}
		case "urban_wotd":
			cmdMap["urban_wotd"] = &CommandModule{nil, HelpUrbanWotd, UrbanWotd}
		case "urban_trending":
			cmdMap["urban_trending"] = &CommandModule{nil, HelpUrbanTrending, UrbanTrending}
		case "url":
			parserModules = append(parserModules, &ParserModule{nil, Url})
		case "interject":
			cmdMap["interject"] = &CommandModule{nil, HelpInterject, Interject}
		case "xkcd":
			cmdMap["xkcd"] = &CommandModule{nil, Helpxkcd, getXKCD}
		case "doing":
			cmdMap["doing"] = &CommandModule{ConfigDoing, HelpDoing, Doing}
		case "doom":
			cmdMap["doom"] = &CommandModule{nil, HelpDoom, Doom}
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
		if m.Configure != nil {
			m.Configure(cfg)
		}
	}

	actions := Actions{
		Quit: func() {
			conn.Quit()
		},
		SayTo: func(dest string, msg string) {
			outboundMsg := message.OutboundMsg{
				Conn: conn,
				Dest: dest,
				Msg:  msg,
			}
			outChan <- outboundMsg
		},
	}

	for _, m := range tickerModules {
		if m.Configure != nil {
			m.Configure(cfg)
		}
		m.Ticker = time.NewTicker(m.GetDuration(cfg))
		tickerChan := m.Ticker.C
		run := m.Run
		go func() {
			// Note that the sender of this channel will never close it.
			// It must be closed manually after time.Stop in order to exit this goroutine.
			for t := range tickerChan {
				run(cfg, t, &actions)
			}
		}()
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

		actions.Say = func(msg string) {
			outboundMsg := message.OutboundMsg{
				Conn: conn,
				Dest: inboundMsg.Src,
				Msg:  msg,
			}
			outChan <- outboundMsg
		}

		// run parser modules
		for _, m := range parserModules {
			if m.Run(cfg, &inboundMsg, &actions) {
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
				m.Run(cfg, &inboundMsg, &actions)
			}
		}

	})

	return nil
}
