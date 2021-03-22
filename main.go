package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	inputPort := flag.String("inputPort", "8000", "port to listen on")
	id := flag.String("id", "no-id", "id for config/log")

	flag.Parse()

	repeaterId := *id

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	os.MkdirAll(dir+"/config", os.ModePerm)
	os.MkdirAll(dir+"/log", os.ModePerm)

	listenners := []Listener{Listener{
		Id:   1,
		Addr: "*",
		Port: *inputPort,
	}}
	transmitters := []Transmitter{Transmitter{
		Id:   2,
		Addr: "*",
		Port: "*",
	}}

	targets := []Target{}
	args := flag.Args()
	for i := range args {
		// ex: 217.122.333.112:5002
		s := strings.Split(args[i], ":")
		targets = append(targets, Target{
			Id:            i + 10,
			Port:          s[1],
			Addr:          s[0],
			TransmitterId: 2,
		})
	}

	var targetIdArray []int
	for i := range targets {
		targetIdArray = append(targetIdArray, targets[i].Id)
	}
	mapForwards := []MapForward{
		MapForward{
			SourceId: 1,
			Addr:     "*",
			Port:     "*",
			Target:   targetIdArray,
		}}

	config := RepeaterConfig{
		Listeners:    listenners,
		Transmitters: transmitters,
		Targets:      targets,
		Maps:         mapForwards,
	}

	c, _ := json.Marshal(config)
	ioutil.WriteFile(dir+"/config/"+repeaterId+".json", c, os.ModePerm)

	repeater := exec.Command(dir+"/repeater", dir+"/config/"+repeaterId+".json", dir+"/log/"+repeaterId+".log")

	repeaterOut, _ := repeater.StdoutPipe()
	repeaterErr, _ := repeater.StderrPipe()

	go func() {
		scanner := bufio.NewScanner(repeaterOut)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(repeaterErr)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()

	err = repeater.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Listener struct {
	Id   int    `json:"id"`
	Addr string `json:"address"`
	Port string `json:"port"`
}

type Transmitter struct {
	Id   int    `json:"id"`
	Addr string `json:"address"`
	Port string `json:"port"`
}

type Target struct {
	Id            int    `json:"id"`
	Addr          string `json:"address"`
	Port          string `json:"port"`
	TransmitterId int    `json:"transmitter"`
}

type MapForward struct {
	SourceId int    `json:"source"`
	Addr     string `json:"address"`
	Port     string `json:"port"`
	Target   []int  `json:"target"`
}

type RepeaterConfig struct {
	Listeners    []Listener    `json:"listen"`
	Transmitters []Transmitter `json:"transmit"`
	Targets      []Target      `json:"target"`
	Maps         []MapForward  `json:"map"`
}
