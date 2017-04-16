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
    cmd_prefix: "!"

connection:
    verbose_callback_handler: False
    debug: False

admin:
    owner: mynick
    deny_message: yeah, i don't think so

# This is a list of modules that irkbot supports. If you omit any of these, they
# will not be loaded at runtime.
modules:
    echo_name:
    help:
    slam:
        # This is the location of the sqlite database used by the slam module.
        db_file: /home/david/var/irkbot/slam.db
    compliment:
        # This is the location of the sqlite database used by the compliment module.
        db_file: /home/david/var/irkbot/compliment.db
    quit:
    quote:
        # This is the location of the sqlite database used by the quotes module.
        db_file: /home/david/var/irkbot/quotes.db
    say:
    urban:
    urban_wotd:
    urban_trending:
    url:
```

### Usage

Once you've created the [configuration file](#configure), simply run the Irkbot binary:

```
$GOPATH/bin/irkbot
```

Once Irkbot has connected, you can get a list of bot commands by typing the help command in either a channel that Irkbot is in or a in a private message. E.g., if the `cmd_prefix` config value is set to `!`, then the help command would be `!help`.

More information about Irkbot's current modules can be gathered by looking at the modules' source in the [lib/module](lib/module) directory.

### Module development

Irkbot has a simple system for creating modules that extend the bot's functionality. Currently, the only modules that Irkbot manages are for PRIVMSG actions, but there are plans to make modules for other events, such as time-based events.

Below is an example PRIVMSG module that adds an echo command to the bot. The echo module runs when a PRIVMSG contains the echo command and sends a PRIVMSG back echoing the rest of the line.

This module file belongs in the [lib/module](lib/module) directory.

```go
// module is the package for all irkbot modules
package module

import (
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

func ConfigEcho(cfg *configure.Config) {
	// This is an optional function to configure the module. It is called only
	// once when irkbot starts up.
	// This function can be omitted if no configuration is needed.
}

func HelpEcho() []string {
	// This function returns an array of strings describing this command's
	// functionality. It is displayed when someone gives the help command in a
	// channel or private message.
	s := "echo <phrase> - echo the phrase back to the channel"
	return []string{s}
}

func Echo(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	// Grab the rest of the message.
	msg := strings.Join(in.MsgArgs[1:], " ")

	// Call the Say function (which does message rate-limiting)
	actions.Say(msg)
}
```

The final step is to add the echo module functions to the switch statement in the RegisterModules function in [lib/module/register.go](lib/module/register.go):

```go
		case "echo":
			cmdMap["echo"] = &CommandModule{ConfigEcho, HelpEcho, Echo}
```

If you omit the config function, the register function call would be:

```go
		case "echo":
			cmdMap["echo"] = &CommandModule{nil, HelpEcho, Echo}
```

To enable the module, you'll need to add it to the `modules` section of the Irkbot configuration file.

### TODO

* Alphabetize help string array.
* Make reddit module.
* Allow multiple channels.
* Add time-based modules.
* Implement max message length handling.
