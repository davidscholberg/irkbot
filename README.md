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
        # These files are used by the "slam" module to pull verbal smackdowns
        # from.
        # The slam module will fail gracefully if these options are missing.
        adjective_file: /home/david/.config/irkbot/slam-adjectives.txt
        noun_file: /home/david/.config/irkbot/slam-nouns.txt
    compliment:
        # This file is used by the "compliment" module to pull compliments from.
        # The compliment module will fail gracefully if this option is missing.
        file: /home/david/.config/irkbot/compliments.txt
    quit:
    quote:
        # This is the location of the sqlite database used by the quotes module.
        db_file: /home/david/var/irkbot/quotes.db
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

Once Irkbot has connected, you can get a list of bot commands by typing `..help` in either a channel that Irkbot is in or a in a private message.

More information about Irkbot's current modules can be gathered by looking at the modules' source in the `lib/modules` directory.

### Module development

Irkbot has a simple system for creating modules that extend the bot's functionality. Currently, the only modules that Irkbot manages are for PRIVMSG actions, but there are plans to make modules for other events, such as time-based events.

Below is an example PRIVMSG module that adds an echo command to the bot. The echo module runs when a PRIVMSG contains the echo command and sends a PRIVMSG back echoing the rest of the line.

This module file belongs in the `lib/module/` directory.

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

func Echo(in *message.InboundMsg, actions *Actions) {
	// Grab the rest of the message.
	msg := strings.Join(in.MsgArgs[1:], " ")

	// Call the Say function (which does message rate-limiting)
	actions.Say(msg)
}
```

The final step is to add the echo module functions to the switch statement in the RegisterModules function in `lib/module/register.go`:

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
