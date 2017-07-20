package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	yaml "gopkg.in/yaml.v2"
)

const (
	filename = "config.yml"
)

var conf config
var refresh = make(chan struct{})

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go handleSIGUSR1(c)
	go reload()

	fmt.Printf("Reload yml file   : kill -SIGUSR1 %s\n", strconv.Itoa(os.Getpid()))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL)
	<-quit
}

type config struct {
	URL string `yaml:"url"`
}

func reload() {
	for range refresh {
		load()
	}
}

func load() error {
	var c config
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	fmt.Println(c)

	conf = c
	return nil
}

func handleSIGUSR1(c chan os.Signal) {
	for {
		<-c
		fmt.Println("got signal SIGUSR1")
		refresh <- struct{}{}
	}
}
