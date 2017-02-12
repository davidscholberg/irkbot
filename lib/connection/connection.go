package connection

import (
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/davidscholberg/irkbot/lib/module"
	"github.com/thoj/go-ircevent"
)

func GetIrcConn(cfg *configure.Config) *irc.Connection {
	conn := irc.IRC(cfg.User.Nick, cfg.User.User)
	conn.UseTLS = cfg.Server.UseTls
	conn.VerboseCallbackHandler = cfg.Connection.VerboseCallbackHandler
	conn.Debug = cfg.Connection.Debug

	conn.AddCallback("001", func(e *irc.Event) {
		if cfg.User.Identify && conn.GetNick() == cfg.User.Nick {
			conn.Privmsgf("NickServ", "identify %s", cfg.User.Password)
		}
		conn.Join(cfg.Channel.ChannelName)
	})

	conn.AddCallback("366", func(e *irc.Event) {
		if len(cfg.Channel.Greeting) != 0 {
			conn.Privmsg(e.Arguments[1], cfg.Channel.Greeting)
		}
	})

	// TODO: start multiple sayLoops, one per conn
	// TODO: pass conn to sayLoop instead of privmsg callbacks?
	sayChan := make(chan message.SayMsg)
	go message.SayLoop(sayChan)

	module.RegisterModules(conn, cfg, sayChan)

	return conn
}
