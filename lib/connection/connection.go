package connection

import (
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/davidscholberg/irkbot/lib/module"
	"github.com/thoj/go-ircevent"
	"time"
)

func GetIrcConn(cfg *configure.Config) (*irc.Connection, error) {
	conn := irc.IRC(cfg.User.Nick, cfg.User.User)
	conn.UseTLS = cfg.Server.UseTls
	conn.VerboseCallbackHandler = cfg.Connection.VerboseCallbackHandler
	conn.Debug = cfg.Connection.Debug

	if cfg.Server.ServerAuth {
		conn.Password = cfg.Server.ServerPassword
	}
	conn.AddCallback("001", func(e *irc.Event) {
		if cfg.User.Identify && conn.GetNick() == cfg.User.Nick {
			conn.Privmsgf("NickServ", "identify %s", cfg.User.Password)
			// temporary horrible hack to allow time to be identified
			// before joining a channel
			time.Sleep(time.Second * 10)
		}
		conn.Join(cfg.Channel.ChannelName)
	})

	conn.AddCallback("366", func(e *irc.Event) {
		if len(cfg.Channel.Greeting) != 0 {
			conn.Privmsg(e.Arguments[1], cfg.Channel.Greeting)
		}
	})

	conn.AddCallback("KICK", func(e *irc.Event) {
		if cfg.Channel.AutoJoinOnKick {
			conn.Join(cfg.Channel.ChannelName)
		}
	})

	outChan := make(chan message.OutboundMsg)
	go message.SayLoop(outChan)

	err := module.RegisterModules(conn, cfg, outChan)
	if err != nil {
		return conn, err
	}

	return conn, nil
}
