package modpm

import (
    "github.com/davidscholberg/irkbot/lib"
)

func RegisterMods(registerMod func(m *lib.Module)) {
    registerMod(&lib.Module{nil, EchoName})
    registerMod(&lib.Module{ConfigInsult, Insult})
    registerMod(&lib.Module{nil, Quit})
    registerMod(&lib.Module{nil, Urban})
}