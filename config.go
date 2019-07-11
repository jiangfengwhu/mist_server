package main

import (
	"encoding/json"
	"log"
	"os"
	"flag"
)

var globalConf configModel

type configModel struct {
	RecapSecure string `json:"recapSecure"`
	Host        string `json:"host"`
	ResDir      string `json:"resourcedir"`
	ResRef      string `json:"resref"`
	Announce    string `json:"announce"`
	RecapServer string `json:"recapServer"`
}

func getConfig() {
	confpath := flag.String("config", "./.config.json", "config file path")
	flag.Parse()
	configFile, err := os.Open(*confpath)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewDecoder(configFile).Decode(&globalConf)
	log.Println(globalConf)
}
