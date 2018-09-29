package module

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"os"
	"strconv"
	"strings"
	"time"
)

type Feature struct {
	ID   uint `gorm:"primary_key"`
	Nick string
	Text string
	Date time.Time
	Live uint
}

func ConfigFeature(cfg *configure.Config) {
	dbFile = cfg.Modules["features"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&Feature{})
}

func HelpRequestFeature() []string {
	s := "featurerequest <string> - add <string> to the list of requested features for an admin to review"
	return []string{s}
}

func RequestFeature(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if len(in.MsgArgs) < 2 {
		actions.Say(
			fmt.Sprintf(
				"%s: you need to include a feature, Einstein",
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

	db.Create(&Feature{Nick: in.Event.Nick, Text: strings.Join(in.MsgArgs[1:], " "), Date: time.Now(), Live: 1})
	actions.Say(fmt.Sprintf("%s: ok, got it", in.Event.Nick))

	return
}

func HelpGetFeatures() []string {
	s := "features - fetch suggested features from the database"
	return []string{s}
}

func GetFeatures(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		actions.Say("couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	features := []Feature{}
	db.Find(&features, Feature{Live: 1})
	if len(features) == 0 {
		actions.Say(
			fmt.Sprintf(
				"%s: no requested features",
				in.Event.Nick,
			),
		)
		return
	}

	if strings.HasPrefix(in.Src, "#") {
		actions.Say(
			fmt.Sprintf(
				"%s: I'll PM you, homie",
				in.Event.Nick,
			),
		)
	}

	for _, feature := range features {
		actions.SayTo(
			in.Event.Nick,
			fmt.Sprintf(
				"%d: %s",
				feature.ID,
				feature.Text,
			),
		)
	}

	return
}

func HelpMarkFeature() []string {
	s := "featuredone <id> - mark feature <id> as done"
	return []string{s}
}

func MarkFeatureDone(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
        if in.Event.Nick != cfg.Admin.Owner {
                actions.Say("%s: who are you, again?", in.Event.Nick)
                return
        }
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		actions.Say("couldn't open quotes database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	featureIdU32, err := strconv.ParseUint(strings.TrimSpace(in.MsgArgs[1]), 10, 32)
	featureId := uint(featureIdU32)
	if err != nil {
		actions.Say("that's not a number, weirdo")
		return
	}

	features := []Feature{}
	db.Find(&features, Feature{ID: featureId})
	if len(features) == 0 {
		actions.Say(
			fmt.Sprintf(
				"%s: no such feature",
				in.Event.Nick,
			),
		)
		return
	}

	if len(features) > 1 {
		actions.Say(
			fmt.Sprintf(
				"%s: you broke something and that returned too many records",
				in.Event.Nick,
			),
		)
		return
	}

	f := Feature{}
	db.First(&f, Feature{ID: featureId})

	if f.Live == 0 {
		actions.Say(
			fmt.Sprintf(
				"%s: that one's already marked as done. Kudos for enthusiasm, though!",
				in.Event.Nick,
			),
		)
		return
	}

	f.Live = 0
	db.Save(&f)

	actions.Say(
		fmt.Sprintf(
			"%s: feature %d marked complete",
			in.Event.Nick,
			featureId,
		),
	)

	return
}
