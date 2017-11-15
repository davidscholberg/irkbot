package module

import (
    "fmt"
    "github.com/jholtom/irkbot/lib/configure"
    "github.com/jholtom/irkbot/lib/message"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "os"
    "strings"
)

type Paste struct {
    ID   uint `gorm:"primary_key"`
    Name string
    Text string
}

var dbPasteFile string

func ConfigPaste(cfg *configure.Config) {
    dbPasteFile = cfg.Modules["paste"]["db_file"]

    db, err := gorm.Open("sqlite3", dbPasteFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer db.Close()

    db.AutoMigrate(&Paste{})
}

func HelpGetPaste() []string {
    s := "paste <name> - Paste the <name>d content from the paste database"
    return []string{s}
}

func GetPaste(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
    if len(in.MsgArgs) < 2 {
        actions.Say(
            fmt.Sprintf(
                "%s: you must specify the <name> of a paste to get",
                in.Event.Nick,
            ),
        )
        return
    }

    db, err := gorm.Open("sqlite3", dbPasteFile)
    if err != nil {
        actions.Say("couldn't open pastes database")
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer db.Close()

    name := strings.TrimSpace(in.MsgArgs[1])

    paste := []Paste{}
    db.Find(&paste, Paste{Name: name})
    if len(paste) == 0 {
        actions.Say(
            fmt.Sprintf(
                "%s: no paste for %s",
                in.Event.Nick,
                name,
            ),
        )
        return
    }

    output := paste[0]
    splat := strings.Split(output.Text, "\n")
    for _, k := range splat {
        msg := fmt.Sprintf(
            "%s",
            k,
        )
        actions.Say(msg)
    }
    return
}
