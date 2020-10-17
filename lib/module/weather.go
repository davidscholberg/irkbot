package module

import (
	"fmt"
	"github.com/briandowns/openweathermap"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

func helpWeather() []string {
	s := "weather <location> - display current weather for the given location (only <city> or <city,country> searches are supported)"
	return []string{s}
}

func weather(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if !strings.HasPrefix(in.Src, "#") {
		actions.say("weather searches not allowed in PMs")
		return
	}

	if len(in.MsgArgs) < 2 {
		actions.say(fmt.Sprintf("%s: please specify a location (<city> or <city,country>)", in.Event.Nick))
		return
	}

	msg := strings.Join(in.MsgArgs[1:], " ")
	//fetch API key from config
	apiKey := cfg.Modules["weather"]["api_key"]

	c := &http.Client{Timeout: time.Duration(cfg.Http.Timeout) * time.Second}
	w, err := openweathermap.NewCurrent(
		"c",
		"en",
		apiKey,
		openweathermap.WithHttpClient(c),
	)
	if err != nil {
		actions.say("error initializing weather search :(")
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = w.CurrentByName(msg)
	if err != nil {
		actions.say("No results returned. Only <city> or <city,country> searches are supported.")
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

	actions.say(
		fmt.Sprintf(
			"current weather for %s, %s: %.0f°C, %d%% humidity, wind %s at %.0fm/s, %s",
			w.Name,
			w.Sys.Country,
			w.Main.Temp,
			w.Main.Humidity,
			degreeToCompassDir(w.Wind.Deg),
			w.Wind.Speed,
			conditions,
		),
	)
}

func degreeToCompassDir(degree float64) string {
	compassDirs := [8]string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	return compassDirs[int(math.Floor((degree+22.5)/45))%8]
}
