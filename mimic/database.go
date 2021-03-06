package mimic

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
)

type Database struct {
	DbName       string   `json:"name"`
	DbHost       string   `json:"host"`
	DbPort       string   `json:"port"`
	DbUser       string   `json:"user"`
	DbExtensions []string `json:"extensions"`
}

func (d Database) Name() string {
	return d.DbName
}

func (d Database) Extensions() []string {
	return d.DbExtensions
}

func (d Database) Host() string {
	return d.DbHost
}

func (d Database) Port() string {
	return d.DbPort
}

func (d Database) User() string {
	if d.DbUser != "" {
		return d.DbUser
	} else {
		username, _ := userName()
		return username
	}
}

func (d Database) Drop() error {
	fmt.Printf("[%s] Dropping database\n", d.Name())
	fmt.Printf(
		"[%v] %v %v %v %v %v %v %v %v\n",
		d.Name(),
		"dropdb",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		d.Name(),
	)
	cmd := exec.Command(
		"dropdb",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		d.Name(),
	)

	return cmdRunner(cmd)
}

func (d Database) Create() error {
	fmt.Printf("[%s] Creating database\n", d.Name())
	fmt.Printf(
		"[%v] %v %v %v %v %v %v %v %v\n",
		d.Name(),
		"createdb",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		d.Name(),
	)
	cmd := exec.Command(
		"createdb",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		d.Name(),
	)

	return cmdRunner(cmd)
}

func (d Database) CreateExtensions() error {
	fmt.Printf("[%s] Creating extensions %v\n", d.Name(), d.Extensions())
	for _, extension := range d.Extensions() {
		extensionCmd := fmt.Sprintf("create extension %s", extension)
		fmt.Printf(
			"[%v] %v %v %v %v %v %v %v %v %v '%v'\n",
			d.Name(),
			"psql",
			"-h",
			d.Host(),
			"-p",
			d.Port(),
			"-U",
			d.User(),
			d.Name(),
			"-c",
			extensionCmd,
		)

		cmd := exec.Command(
			"psql",
			"-h",
			d.Host(),
			"-p",
			d.Port(),
			"-U",
			d.User(),
			d.Name(),
			"-c",
			extensionCmd,
		)

		err := cmdRunner(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

//Restore does a pg_restore to the app from the passed file.
func (d Database) Restore(fileName string) error {
	fmt.Printf("[%s] Restoring database from %v\n", d.Name(), fileName)
	fmt.Printf(
		"[%v] %v %v %v %v %v %v %v %v %v %v %v %v %v %v\n",
		d.Name(),
		"pg_restore",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		"-O",
		"-x",
		"-n",
		"public",
		"-d",
		d.Name(),
		fileName,
	)
	cmd := exec.Command(
		"pg_restore",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		"-O",
		"-x",
		"-n",
		"public",
		"-d",
		d.Name(),
		fileName,
	)

	return cmdRunner(cmd)
}

//Dump takes a pg_dump for app on staging storing it in fileName
func (d Database) Dump(passwd string, fileName string) error {
	fmt.Printf("[%s] Dumping datebase to %v\n", d.Name(), fileName)
	fmt.Printf(
		"[%v] %v %v %v %v %v %v %v %v %v %v %v\n",
		d.Name(),
		"pg_dump",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		"-Fc",
		"-f",
		fileName,
		d.Name(),
	)
	dump := exec.Command(
		"pg_dump",
		"-h",
		d.Host(),
		"-p",
		d.Port(),
		"-U",
		d.User(),
		"-Fc",
		"-f",
		fileName,
		d.Name(),
	)
	dump.Env = append(dump.Env, "PGPASSWORD="+passwd)
	return cmdRunner(dump)
}

func removeFile(fileName string) error {
	cmd := exec.Command("rm", "-f", fileName)
	err := cmdRunner(cmd)
	if err != nil {
		return fmt.Errorf("removeFile: %v", err)
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
