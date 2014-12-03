package icis

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/wm/mimic/application"
)

//App represents a Rails application that is used locally and remotely (on
//staging).
type App struct {
	Name            string   `json:"name"`
	DumpDatabase    Db       `json:"dump_database"`
	RestoreDatabase Db       `json:"restore_database"`
	AfterRakeTasks  []string `json:"after_rake_tasks"`
	DockerApp       bool     `json:"docker_app"`
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

func (a App) DumpDb() application.Database {
	return a.DumpDatabase
}

func (a App) RestoreDb() application.Database {
	return a.RestoreDatabase
}

// the path to the application in development.
func (a App) Path() string {
	homeDir, err := userHomeDir()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get User directory:", err)
		os.Exit(1)
	}

	if a.DockerApp {
		return homeDir + "/src/docker-dev"
	} else {
		return homeDir + "/src/" + a.Name
	}
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

	err = app.RestoreDb().Restore(fileName)
	if err != nil {
		return fmt.Errorf("RestoreDevDatabase: %v", err)
	}

	return err
}

//RunAfterRakeTasks executes the given task in the app
func (a App) RunAfterRakeTasks() error {
	var err error
	var task *exec.Cmd
	var bundle *exec.Cmd

	for _, taskName := range a.AfterRakeTasks {
		if a.DockerApp == true {
			fmt.Printf("[%s] fig run %s bundle exec rake %s in %s\n", a, a, taskName, a.Path())
			bundle = exec.Command("fig", "run", a.Name, "bundle")
			task = exec.Command("fig", "run", a.Name, "bundle", "exec", "rake", taskName)
		} else {
			fmt.Printf("[%s] bundle exec rake %s in %s\n", a, taskName, a.Path())
			bundle = exec.Command("bundle")
			task = exec.Command("bundle", "exec", "rake", taskName)
		}
		bundle.Dir = a.Path()
		task.Dir = a.Path()
		err = cmdRunner(bundle)
		if err != nil {
			return fmt.Errorf("bundle: %v", err)
		}
		err = cmdRunner(task)
		if err != nil {
			return fmt.Errorf("bundle exec rake %v: %v", taskName, err)
		}
	}

	return err
}

func userName() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}

func userHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func cmdRunner(cmd *exec.Cmd) (err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err = cmd.Run(); err != nil {
		return &StderrError{&err, &stderr, &stdout}
	}

	return
}
