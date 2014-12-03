package mimic

import (
	"fmt"
	"os/exec"
)

type Database struct {
	DbName string `json:"name"`
	DbHost string `json:"host"`
	DbPort string `json:"port"`
	DbUser string `json:"user"`
}

func (d Database) Name() string {
	return d.DbName
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
