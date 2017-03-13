package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"math/rand"
	"os"
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

func UpdateQuoteBuffer(p *message.Privmsg) bool {
	// don't update quote buffer in PMs
	if !strings.HasPrefix(p.Dest, "#") {
		return false
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	defer db.Close()

	q := QuoteBuffer{}
	db.FirstOrCreate(&q, QuoteBuffer{Nick: p.Event.Nick})
	db.Model(&q).Updates(QuoteBuffer{Text: p.Msg, Date: time.Now()})

	return false
}

func HelpGrabQuote() []string {
	s := "grab <nick> - store the last quote from <nick> in the quotes database"
	return []string{s}
}

func GrabQuote(p *message.Privmsg) {
	// don't allow quote grabs in PMs
	if !strings.HasPrefix(p.Dest, "#") {
		message.Say(
			p,
			fmt.Sprintf(
				"%s: you can't grab a quote in a PM, doofus",
				p.Event.Nick,
			),
		)
		return
	}

	if len(p.MsgArgs) < 2 {
		message.Say(
			p,
			fmt.Sprintf(
				"%s: you need to specify a nick, dingus",
				p.Event.Nick,
			),
		)
		return
	}

	quotee := strings.TrimSpace(p.MsgArgs[1])

	if p.Event.Nick == quotee {
		message.Say(
			p,
			fmt.Sprintf(
				"%s: you can't grab your own quotes, you narcissistic fool",
				p.Event.Nick,
			),
		)
		return
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		message.Say(p, "couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	q := QuoteBuffer{}
	db.First(&q, QuoteBuffer{Nick: quotee})
	if q.Nick == "" {
		message.Say(
			p,
			fmt.Sprintf(
				"%s: no entries for %s",
				p.Event.Nick,
				quotee,
			),
		)
		return
	}

	db.Create(&Quote{Nick: q.Nick, Text: q.Text, Date: q.Date})
	db.Delete(&q)
	message.Say(p, fmt.Sprintf("%s: grabbed", p.Event.Nick))

	return
}

func HelpGetQuote() []string {
	s := "quote <nick> - get a random quote of <nick> from the quotes database"
	return []string{s}
}

func GetQuote(p *message.Privmsg) {
	if len(p.MsgArgs) < 2 {
		message.Say(
			p,
			fmt.Sprintf(
				"%s: you need to specify a nick, dingus",
				p.Event.Nick,
			),
		)
		return
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		message.Say(p, "couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	quotee := strings.TrimSpace(p.MsgArgs[1])

	quotes := []Quote{}
	db.Find(&quotes, Quote{Nick: quotee})
	if len(quotes) == 0 {
		message.Say(
			p,
			fmt.Sprintf(
				"%s: no entries for %s",
				p.Event.Nick,
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
	message.Say(p, msg)

	return
}
