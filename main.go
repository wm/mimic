package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"syscall"

	"code.google.com/p/gopass"
	"github.com/wm/mimic/mimic"
	"github.com/wm/mimic/utils"
)

type Config struct {
	Applications []mimic.App  `json:"apps"`
	Tunnel       utils.Tunnel `json:"tunnel"`
}

func (cfg *Config) Apps() []mimic.App {
	var apps []mimic.App

	for _, app := range cfg.Applications {
		apps = append(apps, app)
	}

	return apps
}

func (cfg *Config) Load(buffer []byte) {
	err := json.Unmarshal(buffer, cfg)
	if err != nil {
		utils.PrintError("Unmarshal config failed:", err)
		return
	}
}

func main() {
	var cfg Config
	cfgbuf, err := getConfigFile()
	if err != nil {
		utils.PrintError("Read config file failed:", err)
		return
	}
	cfg.Load(cfgbuf)

	selectedApps, err := utils.PromptUserToSelectFrom(cfg.Apps())
	if err != nil {
		utils.PrintError("selecting the applications to update: %v", err)
		return
	}

	stagingPass, err := getStagingPassword()
	if err != nil {
		utils.PrintError("reading your password: %v", err)
		return
	}

	fmt.Printf(
		"[tunnel] tunneling %v to %v:%v through %v\n",
		cfg.Tunnel.LocalPort,
		cfg.Tunnel.RemoteIp,
		cfg.Tunnel.RemotePort,
		cfg.Tunnel.Host,
	)
	err = cfg.Tunnel.Open()

	if err == nil {
		utils.PrintOk("[tunnel] Opened\n")
	} else {
		utils.PrintError("[tunnel] Failed to open tunnel %v\n", err)
	}

	var wg sync.WaitGroup

	for _, app := range selectedApps {
		wg.Add(1)
		go func(app mimic.App) {
			defer wg.Done()
			err := app.RestoreFromDump(stagingPass)
			if err != nil {
				utils.PrintError("[%s] Error restoring! %v\n", app, err)
				return
			}
			utils.PrintOk("[%s] Successfully restored!\n", app)
		}(app)
	}

	wg.Wait()
	err = cfg.Tunnel.Close()

	if err != nil {
		utils.PrintError("[tunnel] Failed to close the tunnel\n")
	}
	utils.PrintOk("[tunnel] Closed\n")
}

func getStagingPassword() (string, error) {
	password, ok := syscall.Getenv("PGPASSWORD_STAGING")

	if ok {
		return password, nil
	}

	return gopass.GetPass("Enter the pg password for staging:")
}

func getConfigFile() ([]byte, error) {
	cfgFile, _ := syscall.Getenv("MIMIC_CONFIG_FILE")
	if len(cfgFile) < 1 {
		cfgFile = "./dev_config.json"
	}

	fmt.Fprintf(os.Stdout, "Enter the configuration file: [%v] ", cfgFile)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(scanner.Text()) > 1 {
		cfgFile = scanner.Text()
	}

	return ioutil.ReadFile(cfgFile)
}
