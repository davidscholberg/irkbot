package modpm

import (
	"github.com/davidscholberg/irkbot/lib"
)

var modules []*lib.Module = []*lib.Module{
	&lib.Module{nil, nil, Help},
	&lib.Module{nil, nil, Url},
	&lib.Module{ConfigEchoName, nil, EchoName},
	&lib.Module{ConfigInsult, HelpInsult, Insult},
	&lib.Module{nil, nil, Quit},
	&lib.Module{nil, HelpUrban, Urban}}

func RegisterMods(registerMod func(m *lib.Module)) {
	for _, m := range modules {
		if m.GetHelp != nil {
			RegisterHelp(m.GetHelp())
		}
		registerMod(m)
	}
}
