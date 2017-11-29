package main

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/connection"
	"log"
	"os"
)

func main() {
	errLogger := log.New(os.Stderr, "error: ", 0)

	// get config
	cfg := configure.Config{}
	err := configure.LoadConfig(&cfg)
	if err != nil {
		errLogger.Fatalln(err)
	}

	conn, err := connection.GetIrcConn(&cfg)
	if err != nil {
		errLogger.Fatalln(err)
	}

	err = conn.Connect(fmt.Sprintf(
		"%s:%s",
		cfg.Server.Host,
		cfg.Server.Port))
	if err != nil {
		errLogger.Fatalln(err)
	}

	conn.Loop()
}
