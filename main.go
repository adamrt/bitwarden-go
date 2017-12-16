package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var allowSignups = flag.Bool("signups", false, "Enable new signups, defaults to false")
var dbFilename = flag.String("db", "", "Path to Database, defaults to ~/.bitwarden.sqlite3")
var port = flag.Int("port", 9996, "Port to use, defaults to 9996")

func main() {
	flag.Parse()

	if *dbFilename == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		*dbFilename = fmt.Sprintf("%s/.bitwarden.sqlite3", user.HomeDir)
	}

	db, err := sqlx.Connect("sqlite3", *dbFilename)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/icons/", iconHandler)                                   // GET
	http.Handle("/identity/connect/token", tokenHandler(db))                  // POST
	http.Handle("/api/accounts/register", registerHandler(db, *allowSignups)) // POST
	http.Handle("/api/sync", syncHandler(db))                                 // GET

	// TODO: Incomplete handlers below
	// r.Post("/api/ciphers", createCipherHandler(db))
	// r.Put("/api/ciphers/:uuid", updateCiphersHandler(db))
	// r.Delete("/api/ciphers/:uuid", deleteCiphersHandler(db))

	// r.Post("/api/folders", createFoldersHandler(db))
	// r.Put("/api/folders/:uuid", updateFoldersHandler(db))
	// r.Delete("/api/folders/:uuid", deleteFoldersHandler(db))

	// r.Put("/api/devices/identifier/:uuid/clear-token", clearTokenHandler(db))
	// r.Put("/api/devices/identifier/:uuid/token", updateTokenHandler(db))

	log.Printf("Listening on 127.0.0.1:%d\n", *port)
	connStr := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(connStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
