package main

import (
	"bufio"
	"bytes"
	"fmt"
	_ "fmt"
	"github.com/spf13/viper"
	"github.com/tmc/scp"
	"log"
	"os"
	"strings"
)

var config = viper.New()
var dv []Device

func main() {
	loadConfig()
	err := prepare()
	if err != nil {
		ExitProgram(err)
	}
	err = scanDevice()
	if err != nil {
		ExitProgram(err)
	}

	for _, k := range dv {
		sshclient, err := connectSSH(k)
		if err != nil {
			fmt.Println("Error connection new SSH session", err)
			_ = sshclient.Close()
			continue
		}
		session, err := sshclient.NewSession()
		if err != nil {
			fmt.Println("NewSession error", err)
			_ = session.Close()
			_ = sshclient.Close()
			continue
		}
		dst := "botinok.rsc"
		err = scp.CopyPath(config.GetString("path.import"), dst, session)
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			fmt.Printf("no such file or directory: %s", dst)
		} else {
			fmt.Println("File copy to Device")
		}
		_ = session.Close()
		sessionN, err := sshclient.NewSession()
		if err != nil {
			fmt.Println("NewSession error", err)
			_ = sessionN.Close()
			_ = sshclient.Close()
			continue
		}
		var b bytes.Buffer
		sessionN.Stdout = &b
		if err := sessionN.Run(config.GetString("ssh.command")); err != nil {
			fmt.Println("Run command error", err)
			_ = sessionN.Close()
			_ = sshclient.Close()
			continue
		}
		if err := sessionN.Run(config.GetString("ssh.command-post")); err != nil {
			fmt.Println("Run command error", err)
			_ = sessionN.Close()
			_ = sshclient.Close()
			continue
		}
		_ = sshclient.Close()
	}
}

func ExitProgram(error error) {
	log.Println(error.Error())
	fmt.Println("Program Exit")
	os.Exit(2)
}
func loadConfig() {
	config.SetConfigFile("./config.yaml")
	config.SetDefault("path.devices", "./router.txt")
	config.SetDefault("path.import", "./botinok.rsc")
	config.SetDefault("path.result", "./result")
	config.SetDefault("path.log", "./log")
	config.SetDefault("ssh.command", "import file=botinok.rsc")
	config.SetDefault("ssh.command-post", "/file remove botinok.rsc")
	config.SetDefault("ssh.port", 22)
	config.SetDefault("ssh.timeout", 30)
	config.SetDefault("delimiter", "::")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//todo
		} else {
			//todo
		}
	}
}
func prepare() error {
	if _, err := os.Stat(config.GetString("path.devices")); os.IsNotExist(err) {
		return errorDeviceFileNotFound
	}
	if _, err := os.Stat(config.GetString("path.import")); os.IsNotExist(err) {
		return errorImportFileNotFound
	}
	if _, err := os.Stat(config.GetString("path.result")); os.IsNotExist(err) {
		err = os.MkdirAll(config.GetString("path.result"), os.ModePerm)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(config.GetString("path.log")); os.IsNotExist(err) {
		err = os.MkdirAll(config.GetString("path.log"), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func scanDevice() error {

	file, err := os.Open(config.GetString("path.devices"))
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // internally, it advances token based on sperator
		s := strings.Split(scanner.Text(), config.GetString("delimiter"))
		if len(s) != 3 {
			return errorStringNotTrueFormat
		}
		x := Device{
			Address:    s[0],
			User:       s[1],
			Password:   s[2],
			SshPort:    uint16(config.GetInt("ssh.port")),
			SshTimeout: config.GetInt("ssh.timeout"),
		}
		dv = append(dv, x)
	}

	return nil
}
