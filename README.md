## Irkbot

Irkbot is a modular IRC bot written in [Go](https://golang.org/) using the [go-ircevent](https://github.com/thoj/go-ircevent) library.

### Get

Fetch and build Irkbot:

```
go get github.com/davidscholberg/irkbot
```

### Configure

Irkbot uses an [INI](https://en.wikipedia.org/wiki/INI_file)-style configuration format. The config file path is expected to be `$HOME/.config/irkbot/irkbot.ini`. Below is a sample configuration file:

```
[user]
nick = mynick
user = mynick

[server]
host = irc.freenode.net
port = 6667

[channel]
channelname = "#mychannel"
greeting = "Sup, folks"

[module]
# This file is used by the "insult" module to pull bad words from.
# The insult module will fail gracefully if this option is missing.
insult-swearfile = "/path/to/badwords.txt"
```

### Usage

Once you've created the [configuration file](#configure), simply run the Irkbot binary:

```
$GOPATH/bin/irkbot
```

Until modules are properly documented, information about Irkbot's current modules can be gathered by looking at the modules' source in the `lib/modules` directory.

### Module development

Irkbot has a simple system for creating modules that extend the bot's functionality. Currently, the only modules that Irkbot manages are for PRIVMSG actions, but there are plans to make modules for other events, such as time-based events.

Below is an example PRIVMSG module that adds an echo command to the bot. The echo module looks for a PRIVMSG beginning with "..echo" and sends a PRIVMSG back echoing the rest of the line.

This module file belongs in the `lib/modules/modpm/` directory.

```golang
// modpm is the package for all PRIVMSG modules
package modpm

import (
    "github.com/davidscholberg/irkbot/lib"
)

func ConfigEcho(cfg *lib.Config) {
    // This is an optional function to configure the module before it is run.
    // This function can be omitted if no configuration is needed.
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

The final step is to add the echo module functions to the modpm.RegisterMods function in `lib/modules/modpm/register.go`:

```golang
    registerMod(&lib.Module{ConfigEcho, Echo})
```

If you omit the config function, the register function call would be:

```golang
    registerMod(&lib.Module{nil, Echo})
```

### TODO

* Remove newlines from HTML title for URL module.
* Make reddit module.
* Allow multiple servers and channels.
* Add time-based modules.
* Add message command to allow messages to be sent to other channels.
* Give each module its own config section.
* Allow each module to specify a command prefix in the config.
* Implement a "help" module that displays the bot's commands.
* Add bot "owner" option and allow  for some modules to be privileged.
* Implement max message length handling.
