## Irkbot

Irkbot is a modular IRC bot written in [Go](https://golang.org/) using the [go-ircevent](https://github.com/thoj/go-ircevent) library. The main goal of this project is to create an IRC bot whose functionality can easily be extended through modules written in pure Go (see [Module development](#module-development)).

### Get

Fetch and build Irkbot:

```
go get github.com/davidscholberg/irkbot
```

### Configure

Irkbot uses [YAML](http://yaml.org/) for its configuration. The config file path is expected to be `$HOME/.config/irkbot/irkbot.yml`. Below is a sample configuration file:

```yaml
user:
    nick: mybot
    user: mybot
    identify: True
    password: mypassword

server:
    host: irc.freenode.net
    port: 7000
    use_tls: True

channel:
    channel_name: "#blahblah"
    greeting: "oh hai"

connection:
    verbose_callback_handler: False
    debug: False

module:
# This file is used by the "insult" module to pull bad words from.
# The insult module will fail gracefully if this option is missing.
    insult_swearfile: /home/david/.config/irkbot/badwords.txt
```

### Usage

Once you've created the [configuration file](#configure), simply run the Irkbot binary:

```
$GOPATH/bin/irkbot
```

Once Irkbot has connected, you can get a list of bot commands by typing `..help` in either a channel that Irkbot is in or a in a private message.

More information about Irkbot's current modules can be gathered by looking at the modules' source in the `lib/modules` directory.

### Module development

Irkbot has a simple system for creating modules that extend the bot's functionality. Currently, the only modules that Irkbot manages are for PRIVMSG actions, but there are plans to make modules for other events, such as time-based events.

Below is an example PRIVMSG module that adds an echo command to the bot. The echo module looks for a PRIVMSG beginning with "..echo" and sends a PRIVMSG back echoing the rest of the line.

This module file belongs in the `lib/modules/modpm/` directory.

```go
// modpm is the package for all PRIVMSG modules
package modpm

import (
	"github.com/davidscholberg/irkbot/lib"
)

func ConfigEcho(cfg *lib.Config) {
	// This is an optional function to configure the module. It is called only
	// once when irkbot starts up.
	// This function can be omitted if no configuration is needed.
}

func HelpEcho() []string {
	// This function returns an array of strings describing this command's
	// functionality. It is displayed when someone types "..help" in a channel
	// or private message.
	s := "..echo <phrase> - echo the phrase back to the channel"
	return []string{s}
}

func Echo(p *lib.Privmsg) bool {
	if ! strings.HasPrefix(p.Msg, "..echo") {
		// If this is not an echo command, return right away.
		// Returning false means that the next module in line will be called.
		return false
	}

	// Grab the rest of the message.
	msg := strings.Join(p.MsgArgs[1:], " ")

	// Call the Say function (which does message rate-limiting)
	lib.Say(p, msg)

	// Returning true causes this module to "consume" this PRIVMSG such that no
	// modules after this one will be called for this PRIVMSG.
	return true
}
```

The final step is to add the echo module functions to the module array in `lib/modules/modpm/register.go`:

```go
	&lib.Module{ConfigEcho, HelpEcho, Echo},
```

If you omit the config function, the register function call would be:

```go
	&lib.Module{nil, HelpEcho, Echo},
```

### TODO

* Implement logging.
* Allow config file to be passed in as an argument.
* Allow each module to specify a command prefix in the config.
* Implement unit testing.
* Allow modules to be disabled in config.
* Alphabetize help string array.
* Make reddit module.
* Allow multiple servers and channels.
* Add time-based modules.
* Add message command to allow messages to be sent to other channels.
* Give each module its own config section.
* Add bot "owner" option and allow for some modules to be privileged.
* Implement max message length handling.
