package module

import (
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/thoj/go-ircevent"
	"strings"
)

type Module struct {
	Configure func(*configure.Config)
	GetHelp   func() []string
	Run       func(*message.Privmsg) bool
}

var modules []*Module = []*Module{
	&Module{nil, nil, Help},
	&Module{nil, nil, Url},
	&Module{ConfigEchoName, nil, EchoName},
	&Module{ConfigInsult, HelpInsult, Insult},
	&Module{nil, nil, Quit},
	&Module{nil, HelpUrban, Urban}}

func RegisterModules(conn *irc.Connection, cfg *configure.Config, sayChan chan message.SayMsg) {
	// register modules
	for _, m := range modules {
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

		for _, m := range modules {
			if m.Run(&p) {
				break
			}
		}
	})

}
