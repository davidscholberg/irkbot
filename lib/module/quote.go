package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type quoteBuffer struct {
	ID   uint   `gorm:"primary_key"`
	Nick string `gorm:"unique_index:idx_quotes_nick"`
	Text string
	Date time.Time
}

type quote struct {
	ID   uint `gorm:"primary_key"`
	Nick string
	Text string
	Date time.Time
}

var dbFile string

func configQuote(cfg *configure.Config) {
	dbFile = cfg.Modules["quote"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&quoteBuffer{})
	db.AutoMigrate(&quote{})
}

func updateQuoteBuffer(cfg *configure.Config, in *message.InboundMsg, actions *actions) bool {
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

	q := quoteBuffer{}
	db.FirstOrCreate(&q, quoteBuffer{Nick: in.Event.Nick})
	db.Model(&q).Updates(quoteBuffer{Text: in.Msg, Date: time.Now()})

	return false
}

func cleanQuoteBuffer(cfg *configure.Config, t time.Time, actions *actions) {
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	err = db.Delete(quoteBuffer{}, "date < datetime('now', '-1 days')").Error
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func getCleanQuoteBufferDuration(cfg *configure.Config) time.Duration {
	secondsStr, ok := cfg.Modules["quote"]["clean_quote_buffer_duration"]
	if !ok {
		secondsStr = "600"
	}
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s, using default time of 600 seconds\n", err)
		seconds = 600
	}
	return time.Second * time.Duration(seconds)
}

func helpGrabQuote() []string {
	s := "grab <nick> - store the last quote from <nick> in the quotes database"
	return []string{s}
}

func grabQuote(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	// don't allow quote grabs in PMs
	if !strings.HasPrefix(in.Src, "#") {
		actions.say(
			fmt.Sprintf(
				"%s: you can't grab a quote in a PM :O",
				in.Event.Nick,
			),
		)
		return
	}

	if len(in.MsgArgs) < 2 {
		actions.say(
			fmt.Sprintf(
				"%s: plz specify a nick",
				in.Event.Nick,
			),
		)
		return
	}

	quotee := strings.TrimSpace(in.MsgArgs[1])

	if in.Event.Nick == quotee {
		actions.say(
			fmt.Sprintf(
				"%s: you can't grab your own quotes :O",
				in.Event.Nick,
			),
		)
		return
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		actions.say("couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	q := quoteBuffer{}
	db.First(&q, quoteBuffer{Nick: quotee})
	if q.Nick == "" {
		actions.say(
			fmt.Sprintf(
				"%s: no entries for %s",
				in.Event.Nick,
				quotee,
			),
		)
		return
	}

	db.Create(&quote{Nick: q.Nick, Text: q.Text, Date: q.Date})
	db.Delete(&q)
	actions.say(fmt.Sprintf("%s: grabbed", in.Event.Nick))

	return
}

func helpGetQuote() []string {
	s := "quote <nick> - get a random quote of <nick> from the quotes database"
	return []string{s}
}

func getQuote(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if len(in.MsgArgs) < 2 {
		actions.say(
			fmt.Sprintf(
				"%s: plz specify a nick",
				in.Event.Nick,
			),
		)
		return
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		actions.say("couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	quotee := strings.TrimSpace(in.MsgArgs[1])

	quotes := []quote{}
	db.Find(&quotes, quote{Nick: quotee})
	if len(quotes) == 0 {
		actions.say(
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
	actions.say(msg)

	return
}
