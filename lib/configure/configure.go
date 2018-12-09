package configure

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	User struct {
		Nick     string `yaml:"nick"`
		User     string `yaml:"user"`
		Identify bool   `yaml:"identify"`
		Password string `yaml:"password"`
	} `yaml:"user"`
	Server struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		UseTls         bool   `yaml:"use_tls"`
		ServerAuth     bool   `yaml:"server_auth"`
		ServerPassword string `yaml:"server_pass"`
	} `yaml:"server"`
	Channel struct {
		ChannelName    string `yaml:"channel_name"`
		Greeting       string `yaml:"greeting"`
		CmdPrefix      string `yaml:"cmd_prefix"`
		AutoJoinOnKick bool   `yaml:"auto_join_on_kick"`
	} `yaml:"channel"`
	Connection struct {
		VerboseCallbackHandler bool `yaml:"verbose_callback_handler"`
		Debug                  bool `yaml:"debug"`
	} `yaml:"connection"`
	Admin struct {
		Owner       string `yaml:"owner"`
		DenyMessage string `yaml:"deny_message"`
	} `yaml:"admin"`
	Http struct {
		ResponseSizeLimit int64  `yaml:"response_size_limit"`
		Timeout           int64  `yaml:"timeout"`
		UserAgent         string `yaml:"user_agent"`
	} `yaml:"http"`
	Modules map[string]map[string]string `yaml:"modules"`
}

const confPathFmt = "%s/.config/irkbot/irkbot.yml"

func LoadConfig(cfg *Config) error {
	confPath := fmt.Sprintf(confPathFmt, os.Getenv("HOME"))
	confStr, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(confStr, cfg)
	if err != nil {
		return err
	}
	return nil
}
