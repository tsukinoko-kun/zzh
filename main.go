package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/zalando/go-keyring"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

const service = "dev.frankmayer.zzh"

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	var c Config

	if os.Args[1] == "set" {
		flag.StringVar(&c.Host, "host", "", "Host")
		flag.IntVar(&c.Port, "port", 22, "Port")
		flag.StringVar(&c.User, "user", "", "User")
		flag.StringVar(&c.Password, "password", "", "Password")
		flag.CommandLine.Parse(os.Args[2:])

		var key string
		if c.Port == 22 {
			key = fmt.Sprintf("%s@%s", c.User, c.Host)
		} else {
			key = fmt.Sprintf("%s@%s:%d", c.User, c.Host, c.Port)
		}
		jsonB, err := json.Marshal(&c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := keyring.Set(service, key, string(jsonB)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Access with key %s saved\n", key)
	} else {
		key := os.Args[1]
		jsonStr, err := keyring.Get(service, key)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := json.Unmarshal([]byte(jsonStr), &c); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cmd := exec.Command("sshpass", "-p", c.Password, "ssh", "-o", "StrictHostKeyChecking=accept-new", "-p", strconv.Itoa(c.Port), fmt.Sprintf("%s@%s", c.User, c.Host))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}
