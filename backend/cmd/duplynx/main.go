package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	dbFile := flag.String("db-file", "var/duplynx.db", "path to SQLite database file")
	addr := flag.String("addr", ":8080", "address for HTTP server to bind")
	embedStatic := flag.Bool("embed-static", true, "serve embedded static assets")
	mode := flag.String("mode", "server", "launch mode: server or seed")

	flag.Parse()

	log.Printf("DupLynx starting (mode=%s, addr=%s, db=%s, embedStatic=%t)\n", *mode, *addr, *dbFile, *embedStatic)

	fmt.Println("DupLynx entrypoint is stubbed; implementation coming in later phases.")
}
