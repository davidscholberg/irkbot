package modpm

import (
	"github.com/davidscholberg/irkbot/lib"
)

func RegisterMods(registerMod func(m *lib.Module)) {
	modules := []*lib.Module{
		&lib.Module{nil, Url},
		&lib.Module{ConfigEchoName, EchoName},
		&lib.Module{ConfigInsult, Insult},
		&lib.Module{nil, Quit},
		&lib.Module{nil, Urban}}

	for _, m := range modules {
		registerMod(m)
	}
}
