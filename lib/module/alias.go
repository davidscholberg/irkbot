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

type Alias struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique_index:idx_aliases_name"`
	Text string
	Date time.Time
}

func ConfigAlias(cfg *configure.Config) {
	dbFile := cfg.Modules["alias"]["db_file"]

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&Alias{})
}

func HelpCreateAlias() []string {
	s := "createalias <alias-name> <alias-text> - create an alias, which will be treated as a command that returns the alias text"
	return []string{s}
}

func HelpDeleteAlias() []string {
	s := "deletealias <alias-name> - delete an alias"
	return []string{s}
}

func HelpListAliases() []string {
	s := "listaliases - list all aliases"
	return []string{s}
}

func GetAliasDB(cfg *configure.Config) (*gorm.DB, error) {
	dbFile := cfg.Modules["alias"]["db_file"]
	return gorm.Open("sqlite3", dbFile)
}

func CheckAliases(cfg *configure.Config, in *message.InboundMsg, actions *Actions) bool {
	cmdPrefix := cfg.Channel.CmdPrefix
	if cmdPrefix == "" {
		cmdPrefix = "."
	}
	if !strings.HasPrefix(in.Msg, cmdPrefix) {
		return false
	}
	aliasName := strings.TrimPrefix(in.MsgArgs[0], cmdPrefix)

	db, err := GetAliasDB(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	defer db.Close()

	alias := Alias{}
	db.First(&alias, Alias{Name: aliasName})
	if len(alias.Text) == 0 {
		return false
	}

	actions.Say(alias.Text)

	return false
}

func CreateAlias(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}

	if len(in.MsgArgs) < 3 {
		actions.Say(fmt.Sprintf("%s: usage: <alias-name> <alias-text>", in.Event.Nick))
		return
	}

	aliasName := strings.TrimSpace(in.MsgArgs[1])
	aliasText := strings.Join(in.MsgArgs[2:], " ")

	db, err := GetAliasDB(cfg)
	if err != nil {
		actions.Say("couldn't open alias database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	alias := Alias{}
	db.FirstOrCreate(&alias, Alias{Name: aliasName})
	if len(alias.Text) != 0 {
		actions.Say(fmt.Sprintf("%s: alias \"%s\" already exists", in.Event.Nick, aliasName))
		return
	}
	db.Model(&alias).Updates(Alias{Text: aliasText, Date: time.Now()})

	actions.Say(fmt.Sprintf("%s: alias \"%s\" created", in.Event.Nick, aliasName))
}

func DeleteAlias(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}

	if len(in.MsgArgs) < 2 {
		actions.Say(fmt.Sprintf("%s: usage: <alias-name>", in.Event.Nick))
		return
	}

	aliasName := strings.TrimSpace(in.MsgArgs[1])

	db, err := GetAliasDB(cfg)
	if err != nil {
		actions.Say("couldn't open alias database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	alias := Alias{}
	db.First(&alias, Alias{Name: aliasName})
	if len(alias.Text) == 0 {
		actions.Say(fmt.Sprintf("%s: alias \"%s\" doesn't exist", in.Event.Nick, aliasName))
		return
	}

	db.Delete(&alias)
	actions.Say(fmt.Sprintf("%s: alias \"%s\" deleted", in.Event.Nick, aliasName))
}

func ListAliases(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	db, err := GetAliasDB(cfg)
	if err != nil {
		actions.Say("couldn't open alias database")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	aliases := []Alias{}
	db.Select("name").Find(&aliases)

	if len(aliases) == 0 {
		actions.Say(fmt.Sprintf("%s: no aliases found", in.Event.Nick))
		return
	}

	aliasesStrings := []string{}
	for _, alias := range aliases {
		aliasesStrings = append(aliasesStrings, alias.Name)
	}
	aliasesString := strings.Join(aliasesStrings, ", ")

	actions.Say(fmt.Sprintf("%s: current list of aliases: %s", in.Event.Nick, aliasesString))
}
