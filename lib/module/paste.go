package module

import (
    "fmt"
    "github.com/davidscholberg/irkbot/lib/configure"
    "github.com/davidscholberg/irkbot/lib/message"
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

var dbFile string

func ConfigPaste(cfg *configure.Config) {
    dbFile = cfg.Modules["paste"]["db_file"]

    db, err := gorm.Open("sqlite3", dbFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer db.Close()

    db.AutoMigrate(&Paste{})
}

func HelpStorePaste() []string {
    s := "store <name> <content> - store the <content> as <name> in the pastes database"
    return []string{s}
}

func HelpGetPaste() []string {
    s := "paste <name> - Paste the <name>d content from the paste database"
    return []string{s}
}

func StorePaste(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
    if len(in.MsgArgs) < 3 {
        actions.Say(
            fmt.Sprintf(
                "%s: you must specify the <name> and <content> of a paste to store",
                in.Event.Nick,
            ),
        )
        return
    }
    name := strings.TrimSpace(in.MsgArgs[1])
    content := strings.TrimSpace(in.MsgArgs[1])

    db, err := gorm.Open("sqlite3", dbFile)
    if err != nil {
        actions.Say("couldn't open paste database")
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer db.Close()
    search := []Paste{}
    db.Find(&search, Paste{Name: name})
    if len(search) != 0 {
        actions.Say(
            fmt.Sprintf(
                "%s: Paste %s already exists.",
                in.Event.Nick,
                name,
            ),
        )
        return
    }
    db.Create(&Paste{Name: name, Text: content})
    actions.Say(fmt.Sprintf("%s: stored paste %s", in.Event.Nick, name))
    return
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

    db, err := gorm.Open("sqlite3", dbFile)
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
    msg := fmt.Sprintf(
        "%s",
        output.Text,
    )
    actions.Say(msg)

    return
}
