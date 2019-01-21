package main

import (
	"encoding/json"
	"log"
	"os"
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
	configFile, err := os.Open("./.config.json")
	if err != nil {
		log.Println(err)
		return
	}
	json.NewDecoder(configFile).Decode(&globalConf)
	log.Println(globalConf)
}
