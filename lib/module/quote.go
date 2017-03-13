package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"time"
)

type QuoteBuffer struct {
	gorm.Model
	Nick string `gorm:"unique_index:idx_quotes_nick"`
	Text string
	Date time.Time
}

var dbFile string

func ConfigQuoteBuffer(cfg *configure.Config) {
	dbFile = cfg.Modules["quote"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&QuoteBuffer{})
}

func UpdateQuoteBuffer(p *message.Privmsg) bool {
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
