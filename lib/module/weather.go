package module

import (
	"fmt"
	"github.com/briandowns/openweathermap"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"os"
	"strings"
)

func HelpWeather() []string {
	s := "weather <location> - display current weather for the given location (only <city> or <city,country> searches are supported)"
	return []string{s}
}

func Weather(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if !strings.HasPrefix(in.Src, "#") {
		actions.Say("weather searches not allowed in PMs")
		return
	}

	if len(in.MsgArgs) < 2 {
		actions.Say(fmt.Sprintf("%s: please specify a location (<city> or <city,country>)", in.Event.Nick))
		return
	}

	msg := strings.Join(in.MsgArgs[1:], " ")
	//fetch API key from config
	apiKey := cfg.Modules["weather"]["api_key"]

	w, err := openweathermap.NewCurrent("c", "en", apiKey)
	if err != nil {
		actions.Say("error initializing weather search :(")
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = w.CurrentByName(msg)
	if err != nil {
		actions.Say("No results returned. Only <city> or <city,country> searches are supported.")
		fmt.Fprintln(os.Stderr, err)
		return
	}

	conditions := ""
	for i, condition := range w.Weather {
		if i > 0 {
			conditions += ", "
		}
		conditions += condition.Description
	}

	actions.Say(
		fmt.Sprintf(
			"current weather for %s, %s: %.2f°C, %d%% humidity, %s",
			w.Name,
			w.Sys.Country,
			w.Main.Temp,
			w.Main.Humidity,
			conditions,
		),
	)
}
