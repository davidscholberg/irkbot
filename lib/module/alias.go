package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"strings"
	"time"
)

type alias struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique_index:idx_aliases_name"`
	Text string
	Date time.Time
}

func configAlias(cfg *configure.Config) {
	dbFile := cfg.Modules["alias"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&alias{})
}

func helpCreateAlias() []string {
	s := "createalias <alias-name> <alias-text> - create an alias, which will be treated as a command that returns the alias text"
	return []string{s}
}

func helpDeleteAlias() []string {
	s := "deletealias <alias-name> - delete an alias"
	return []string{s}
}

func helpListAliases() []string {
	s := "listaliases - list all aliases"
	return []string{s}
}

func getAliasDB(cfg *configure.Config) (*gorm.DB, error) {
	dbFile := cfg.Modules["alias"]["db_file"]
	return gorm.Open("sqlite3", dbFile)
}

func checkAliases(cfg *configure.Config, in *message.InboundMsg, actions *actions) bool {
	cmdPrefix := cfg.Channel.CmdPrefix
	if cmdPrefix == "" {
		cmdPrefix = "."
	}
	if !strings.HasPrefix(in.Msg, cmdPrefix) {
		return false
	}
	aliasName := strings.TrimPrefix(in.MsgArgs[0], cmdPrefix)

	db, err := getAliasDB(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	defer db.Close()

	aliasRow := alias{}
	db.First(&aliasRow, alias{Name: aliasName})
	if len(aliasRow.Text) == 0 {
		return false
	}

	actions.say(aliasRow.Text)

	return false
}

func createAlias(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}

	if len(in.MsgArgs) < 3 {
		actions.say(fmt.Sprintf("%s: usage: <alias-name> <alias-text>", in.Event.Nick))
		return
	}

	aliasName := strings.TrimSpace(in.MsgArgs[1])
	aliasText := strings.Join(in.MsgArgs[2:], " ")

	db, err := getAliasDB(cfg)
	if err != nil {
		actions.say("couldn't open alias database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	aliasRow := alias{}
	db.FirstOrCreate(&aliasRow, alias{Name: aliasName})
	if len(aliasRow.Text) != 0 {
		actions.say(fmt.Sprintf("%s: alias \"%s\" already exists", in.Event.Nick, aliasName))
		return
	}
	db.Model(&aliasRow).Updates(alias{Text: aliasText, Date: time.Now()})

	actions.say(fmt.Sprintf("%s: alias \"%s\" created", in.Event.Nick, aliasName))
}

func deleteAlias(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}

	if len(in.MsgArgs) < 2 {
		actions.say(fmt.Sprintf("%s: usage: <alias-name>", in.Event.Nick))
		return
	}

	aliasName := strings.TrimSpace(in.MsgArgs[1])

	db, err := getAliasDB(cfg)
	if err != nil {
		actions.say("couldn't open alias database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	aliasRow := alias{}
	db.First(&aliasRow, alias{Name: aliasName})
	if len(aliasRow.Text) == 0 {
		actions.say(fmt.Sprintf("%s: alias \"%s\" doesn't exist", in.Event.Nick, aliasName))
		return
	}

	db.Delete(&aliasRow)
	actions.say(fmt.Sprintf("%s: alias \"%s\" deleted", in.Event.Nick, aliasName))
}

func listAliases(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	db, err := getAliasDB(cfg)
	if err != nil {
		actions.say("couldn't open alias database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	aliases := []alias{}
	db.Select("name").Find(&aliases)

	if len(aliases) == 0 {
		actions.say(fmt.Sprintf("%s: no aliases found", in.Event.Nick))
		return
	}

	aliasesStrings := []string{}
	for _, alias := range aliases {
		aliasesStrings = append(aliasesStrings, alias.Name)
	}
	aliasesString := strings.Join(aliasesStrings, ", ")

	actions.say(fmt.Sprintf("%s: current list of aliases: %s", in.Event.Nick, aliasesString))
}
