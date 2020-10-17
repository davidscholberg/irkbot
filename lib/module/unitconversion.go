package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/martinlindhe/unit"
	"os"
	"strconv"
	"strings"
)

func helpC2F() []string {
	s := "c2f <temperature> - convert temperature in celsius to fahrenheit"
	return []string{s}
}

func helpF2C() []string {
	s := "f2c <temperature> - convert temperature in fahrenheit to celsius"
	return []string{s}
}

func c2F(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if len(in.MsgArgs) < 2 {
		actions.say(fmt.Sprintf("%s: please specify a temperature in celsius", in.Event.Nick))
		return
	}

	tempCelsius, err := strconv.ParseFloat(strings.TrimRight(in.MsgArgs[1], "°cC"), 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing string to float: %s\n", err)
		actions.say("error parsing input to float")
		return
	}

	c := unit.FromCelsius(tempCelsius)
	actions.say(fmt.Sprintf("%.0f°F", c.Fahrenheit()))
}

func f2C(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if len(in.MsgArgs) < 2 {
		actions.say(fmt.Sprintf("%s: please specify a temperature in fahrenheit", in.Event.Nick))
		return
	}

	tempFahrenheit, err := strconv.ParseFloat(strings.TrimRight(in.MsgArgs[1], "°fF"), 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing string to float: %s\n", err)
		actions.say("error parsing input to float")
		return
	}

	f := unit.FromFahrenheit(tempFahrenheit)
	actions.say(fmt.Sprintf("%.0f°C", f.Celsius()))
}
