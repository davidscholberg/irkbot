package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/thoj/go-ircevent"
	"io"
	"net/http"
	"strings"
	"time"
)

type commandModule struct {
	configure func(*configure.Config)
	getHelp   func() []string
	run       func(*configure.Config, *message.InboundMsg, *actions)
}

type parserModule struct {
	configure func(*configure.Config)
	run       func(*configure.Config, *message.InboundMsg, *actions) bool
}

type tickerModule struct {
	configure   func(*configure.Config)
	getDuration func(*configure.Config) time.Duration
	run         func(*configure.Config, time.Time, *actions)
	ticker      *time.Ticker
}

type actions struct {
	httpGet  func(string) (*http.Response, error)
	httpPost func(string, string, io.Reader) (*http.Response, error)
	quit     func()
	say      func(string)
	sayTo    func(string, string)
}

func RegisterModules(conn *irc.Connection, cfg *configure.Config, outChan chan message.OutboundMsg) error {
	cmdMap := make(map[string]*commandModule)
	parserModules := []*parserModule{}
	tickerModules := []*tickerModule{}
	for moduleName, _ := range cfg.Modules {
		switch moduleName {
		case "alias":
			cmdMap["createalias"] = &commandModule{nil, helpCreateAlias, createAlias}
			cmdMap["deletealias"] = &commandModule{nil, helpDeleteAlias, deleteAlias}
			cmdMap["listaliases"] = &commandModule{nil, helpListAliases, listAliases}
			parserModules = append(parserModules, &parserModule{configAlias, checkAliases})
		case "direct_message_log":
			parserModules = append(parserModules, &parserModule{nil, directMessageLog})
		case "echo_name":
			parserModules = append(parserModules, &parserModule{nil, echoName})
		case "help":
			cmdMap["help"] = &commandModule{nil, nil, help}
			parserModules = append(parserModules, &parserModule{nil, parseHelp})
		case "slam":
			cmdMap["slam"] = &commandModule{configSlam, helpSlam, slam}
		case "compliment":
			cmdMap["compliment"] = &commandModule{configCompliment, helpCompliment, giveCompliment}
		case "quit":
			cmdMap["quit"] = &commandModule{nil, helpQuit, quit}
		case "quote":
			cmdMap["grab"] = &commandModule{nil, helpGrabQuote, grabQuote}
			cmdMap["quote"] = &commandModule{nil, helpGetQuote, getQuote}
			parserModules = append(parserModules, &parserModule{configQuote, updateQuoteBuffer})
			tickerModules = append(
				tickerModules,
				&tickerModule{
					nil,
					getCleanQuoteBufferDuration,
					cleanQuoteBuffer,
					nil,
				},
			)
		case "say":
			cmdMap["say"] = &commandModule{nil, helpSay, say}
		case "urban":
			cmdMap["urban"] = &commandModule{nil, helpUrban, urban}
		case "urban_wotd":
			cmdMap["urban_wotd"] = &commandModule{nil, helpUrbanWotd, urbanWotd}
		case "urban_trending":
			cmdMap["urban_trending"] = &commandModule{nil, helpUrbanTrending, urbanTrending}
		case "url":
			parserModules = append(parserModules, &parserModule{nil, parseUrls})
		case "interject":
			cmdMap["interject"] = &commandModule{nil, helpInterject, interject}
		case "xkcd":
			cmdMap["xkcd"] = &commandModule{nil, helpxkcd, getXKCD}
		case "doing":
			cmdMap["doing"] = &commandModule{configDoing, helpDoing, doing}
		case "doom":
			cmdMap["doom"] = &commandModule{nil, helpDoom, doom}
		case "unit":
			cmdMap["c2f"] = &commandModule{nil, helpC2F, c2F}
			cmdMap["f2c"] = &commandModule{nil, helpF2C, f2C}
		case "weather":
			cmdMap["weather"] = &commandModule{nil, helpWeather, weather}
		case "youtube":
			cmdMap["yt"] = &commandModule{nil, helpYoutubeSearch, youtubeSearch}
		default:
			return fmt.Errorf("invalid name '%s' in module config", moduleName)
		}
	}

	for _, m := range cmdMap {
		if m.getHelp != nil {
			registerHelp(m.getHelp())
		}
		if m.configure != nil {
			m.configure(cfg)
		}
	}
	sortHelp()

	for _, m := range parserModules {
		if m.configure != nil {
			m.configure(cfg)
		}
	}

	// global http client
	httpClient := &http.Client{Timeout: time.Duration(cfg.Http.Timeout) * time.Second}

	// global http request function
	httpRequest := func(method string, url string, contentType string, body io.Reader) (*http.Response, error) {
		request, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		if cfg.Http.UserAgent != "" {
			request.Header.Set("User-Agent", cfg.Http.UserAgent)
		}
		if contentType != "" {
			request.Header.Set("Content-Type", contentType)
		}
		return httpClient.Do(request)
	}

	actions := actions{
		httpGet: func(url string) (*http.Response, error) {
			return httpRequest("GET", url, "", nil)
		},
		httpPost: func(url string, contentType string, body io.Reader) (*http.Response, error) {
			return httpRequest("POST", url, contentType, body)
		},
		quit: func() {
			conn.Quit()
		},
		sayTo: func(dest string, msg string) {
			outboundMsg := message.OutboundMsg{
				Conn: conn,
				Dest: dest,
				Msg:  msg,
			}
			outChan <- outboundMsg
		},
	}

	for _, m := range tickerModules {
		if m.configure != nil {
			m.configure(cfg)
		}
		m.ticker = time.NewTicker(m.getDuration(cfg))
		tickerChan := m.ticker.C
		run := m.run
		go func() {
			// Note that the sender of this channel will never close it.
			// It must be closed manually after time.Stop in order to exit this goroutine.
			for t := range tickerChan {
				run(cfg, t, &actions)
			}
		}()
	}

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		// check ignore list
		for _, user := range cfg.Ignore.UsersToIgnore {
			if e.Nick == user {
				return
			}
		}

		inboundMsg := message.InboundMsg{}
		inboundMsg.Msg = e.Message()
		inboundMsg.MsgArgs = strings.Fields(inboundMsg.Msg)
		inboundMsg.Src = e.Arguments[0]
		if !strings.HasPrefix(inboundMsg.Src, "#") {
			inboundMsg.Src = e.Nick
		}
		inboundMsg.Event = e

		actions.say = func(msg string) {
			outboundMsg := message.OutboundMsg{
				Conn: conn,
				Dest: inboundMsg.Src,
				Msg:  msg,
			}
			outChan <- outboundMsg
		}

		// run parser modules
		for _, m := range parserModules {
			if m.run(cfg, &inboundMsg, &actions) {
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
				m.run(cfg, &inboundMsg, &actions)
			}
		}

	})

	return nil
}
