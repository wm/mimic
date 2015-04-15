package mimic

import (
	"bytes"
	"fmt"
	"time"
)

//App represents a Rails application that is used locally and remotely (on
//staging).
type App struct {
	Name            string   `json:"name"`
	DumpDatabase    Database `json:"dump_database"`
	RestoreDatabase Database `json:"restore_database"`
}

type StderrError struct {
	BaseError *error
	Stderr    *bytes.Buffer
	Stdout    *bytes.Buffer
}

func (se StderrError) Error() string {
	return fmt.Sprint(se.BaseError) +
		": " + se.Stderr.String() + ": " + se.Stdout.String()
}

func (a App) String() string {
	return a.Name
}

func (a App) DumpDb() Database {
	return a.DumpDatabase
}

func (a App) RestoreDb() Database {
	return a.RestoreDatabase
}

//RestoreFromDump dumps the staging db for app and recreates the local
//app database provisioned with the dump. It then migrates the database.
func (app App) RestoreFromDump(passwd string) error {
	fileName := "/tmp/" + app.String() + "-" + time.Now().Format("2006-01-13_0304")

	err := app.DumpDb().Dump(passwd, fileName)
	if err != nil {
		return fmt.Errorf("Dump: %v", err)
	}
	defer removeFile(fileName)

	err = app.RestoreDb().Drop()
	if err != nil {
		fmt.Printf("[%s] dropdb: %v\n", app, err)
	}

	err = app.RestoreDb().Create()
	if err != nil {
		fmt.Printf("[%s] createdb: %v\n", app, err)
	}

	err = app.RestoreDb().CreateExtensions()
	if err != nil {
		fmt.Printf("[%s] create extension: %v\n", app, err)
	}

	err = app.RestoreDb().Restore(fileName)
	if err != nil {
		return fmt.Errorf("RestoreDevDatabase: %v", err)
	}

	return err
}
