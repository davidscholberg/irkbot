package module

import (
	"fmt"
	"github.com/dvdmuckle/irkbot/lib/configure"
	"github.com/dvdmuckle/irkbot/lib/message"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Adjective struct {
	ID   uint   `gorm:"primary_key"`
	Word string `gorm:"unique_index:idx_adjective_word"`
}

type Noun struct {
	ID   uint   `gorm:"primary_key"`
	Word string `gorm:"unique_index:idx_noun_word"`
}

func ConfigSlam(cfg *configure.Config) {
	dbFile := cfg.Modules["slam"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&Adjective{})
	db.AutoMigrate(&Noun{})

	// seed rng
	rand.Seed(time.Now().UnixNano())
}

func HelpSlam() []string {
	s := "slam [victim] - give the victim a verbal smackdown (or self if no victim" +
		" specified)"
	return []string{s}
}

func Slam(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	dbFile := cfg.Modules["slam"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.Say("error: couldn't open slam database")
		return
	}
	defer db.Close()

	adjectives := []Adjective{}
	nouns := []Noun{}
	db.Find(&adjectives)
	db.Find(&nouns)

	if len(adjectives) == 0 || len(nouns) == 0 {
		actions.Say("error: no smackdowns found")
		return
	}

	victim := in.Event.Nick
	if len(in.MsgArgs) > 1 {
		victim = strings.TrimSpace(strings.Join(in.MsgArgs[1:], " "))
	}

	response := fmt.Sprintf(
		"%s: u %s %s",
		victim,
		adjectives[rand.Intn(len(adjectives))].Word,
		nouns[rand.Intn(len(nouns))].Word)

	actions.Say(response)
}
