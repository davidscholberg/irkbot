package lib

import (
	goirc "github.com/thoj/go-ircevent"
)

// TODO: give modules their own sections in config?
type Config struct {
	User struct {
		Nick string `yaml:"nick"`
		User string `yaml:"user"`
	} `yaml:"user"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Channel struct {
		ChannelName string `yaml:"channel_name"`
		Greeting    string `yaml:"greeting"`
	} `yaml:"channel"`
	Connection struct {
		VerboseCallbackHandler bool `yaml:"verbose_callback_handler"`
		Debug                  bool `yaml:"debug"`
	} `yaml:"connection"`
	Module struct {
		InsultSwearfile string `yaml:"insult_swearfile"`
	} `yaml:"module"`
}

type Module struct {
	Configure func(*Config)
	GetHelp   func() []string
	Run       func(*Privmsg) bool
}

type Privmsg struct {
	Msg     string
	MsgArgs []string
	Dest    string
	Event   *goirc.Event
	Conn    *goirc.Connection
	SayChan chan SayMsg
}

type SayMsg struct {
	Conn *goirc.Connection
	Dest string
	Msg  string
}
