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

type Compliment struct {
	ID   uint   `gorm:"primary_key"`
	Text string `gorm:"unique_index:idx_compliment_text"`
}

var compliments []Compliment

func ConfigCompliment(cfg *configure.Config) {
	dbFile := cfg.Modules["compliment"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&Compliment{})

	db.Find(&compliments)

	// seed rng
	rand.Seed(time.Now().UnixNano())
}

func HelpCompliment() []string {
	s := "compliment [recipient] - give the recipient a compliment (or self if no recipient" +
		" specified)"
	return []string{s}
}

func GiveCompliment(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if len(compliments) == 0 {
		actions.Say("error: no compliments loaded")
		return
	}

	recipient := in.Event.Nick
	if len(in.MsgArgs) > 1 {
		recipient = strings.TrimSpace(strings.Join(in.MsgArgs[1:], " "))
	}

	response := fmt.Sprintf(
		"%s: %s",
		recipient,
		compliments[rand.Intn(len(compliments))].Text,
	)

	actions.Say(response)
}
