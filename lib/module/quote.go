package module

import (
	"fmt"
	"github.com/jholtom/irkbot/lib/configure"
	"github.com/jholtom/irkbot/lib/message"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type QuoteBuffer struct {
	ID   uint   `gorm:"primary_key"`
	Nick string `gorm:"unique_index:idx_quotes_nick"`
	Text string
	Date time.Time
}

type Quote struct {
	ID   uint `gorm:"primary_key"`
	Nick string
	Text string
	Date time.Time
}

var dbFile string

func ConfigQuote(cfg *configure.Config) {
	dbFile = cfg.Modules["quote"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&QuoteBuffer{})
	db.AutoMigrate(&Quote{})
}

func UpdateQuoteBuffer(cfg *configure.Config, in *message.InboundMsg, actions *Actions) bool {
	// don't update quote buffer in PMs
	if !strings.HasPrefix(in.Src, "#") {
		return false
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	defer db.Close()

	q := QuoteBuffer{}
	db.FirstOrCreate(&q, QuoteBuffer{Nick: in.Event.Nick})
	db.Model(&q).Updates(QuoteBuffer{Text: in.Msg, Date: time.Now()})

	return false
}

func CleanQuoteBuffer(cfg *configure.Config, t time.Time, actions *Actions) {
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	err = db.Delete(QuoteBuffer{}, "date < datetime('now', '-1 days')").Error
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func GetCleanQuoteBufferDuration(cfg *configure.Config) time.Duration {
	secondsStr, ok := cfg.Modules["quote"]["clean_quote_buffer_duration"]
	if !ok {
		secondsStr = "600"
	}
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s, using default time of 600 seconds", err)
		seconds = 600
	}
	return time.Second * time.Duration(seconds)
}

func HelpGrabQuote() []string {
	s := "grab <nick> - store the last quote from <nick> in the quotes database"
	return []string{s}
}

func GrabQuote(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	// don't allow quote grabs in PMs
	if !strings.HasPrefix(in.Src, "#") {
		actions.Say(
			fmt.Sprintf(
				"%s: you can't grab a quote in a PM, doofus",
				in.Event.Nick,
			),
		)
		return
	}

	if len(in.MsgArgs) < 2 {
		actions.Say(
			fmt.Sprintf(
				"%s: you need to specify a nick, dingus",
				in.Event.Nick,
			),
		)
		return
	}

	quotee := strings.TrimSpace(in.MsgArgs[1])

	if in.Event.Nick == quotee {
		actions.Say(
			fmt.Sprintf(
				"%s: you can't grab your own quotes, you narcissistic fool",
				in.Event.Nick,
			),
		)
		return
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		actions.Say("couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	q := QuoteBuffer{}
	db.First(&q, QuoteBuffer{Nick: quotee})
	if q.Nick == "" {
		actions.Say(
			fmt.Sprintf(
				"%s: no entries for %s",
				in.Event.Nick,
				quotee,
			),
		)
		return
	}

	db.Create(&Quote{Nick: q.Nick, Text: q.Text, Date: q.Date})
	db.Delete(&q)
	actions.Say(fmt.Sprintf("%s: grabbed", in.Event.Nick))

	return
}

func HelpGetQuote() []string {
	s := "quote <nick> - get a random quote of <nick> from the quotes database"
	return []string{s}
}

func GetQuote(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if len(in.MsgArgs) < 2 {
		actions.Say(
			fmt.Sprintf(
				"%s: you need to specify a nick, dingus",
				in.Event.Nick,
			),
		)
		return
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		actions.Say("couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	quotee := strings.TrimSpace(in.MsgArgs[1])

	quotes := []Quote{}
	db.Find(&quotes, Quote{Nick: quotee})
	if len(quotes) == 0 {
		actions.Say(
			fmt.Sprintf(
				"%s: no entries for %s",
				in.Event.Nick,
				quotee,
			),
		)
		return
	}

	quote := quotes[rand.Intn(len(quotes))]
	msg := fmt.Sprintf(
		"<%s> %s",
		quote.Nick,
		quote.Text,
	)
	actions.Say(msg)

	return
}
