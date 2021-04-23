package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.com/voip-services/go-kamailio-api/api/controllers"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

var app controllers.App

func TestMain(m *testing.M) {

	dbURL := fmt.Sprintf("%s?sslmode=disable", os.Getenv("KAM_DB_URL"))

	//log.Printf("db %s", dbUrl)

	app.Initialize(dbURL)

	app.NewClient("go-kamailio-api_devcontainer_mmock_1:8082") // init mmock

	mgr, err := migrate.New(
		"file://migrates",
		dbURL)
	if err != nil {
		log.Fatal(err)
	}

	migrateUp(mgr)

	code := m.Run()

	migrateDown(mgr)
	os.Exit(code)
}

func migrateUp(mgr *migrate.Migrate) {
	// migration up
	if err := mgr.Up(); err != nil {
		log.Println(err)
	}
}

func migrateDown(mgr *migrate.Migrate) {
	// migrations down
	if err := mgr.Down(); err != nil {
		log.Println(err)
	}

}

// TODO: http test
