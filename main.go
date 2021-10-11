package main

import (
	"log"
	"net/http"
)

type PlayerInfo struct {
	Nickname   string
	PlayerID   string
	Elo        int
	Avatar     string
	SkillLevel string
}

type Player struct {
	PlayerID string `json:"player_id"`
	Games    Games  `json:"games"`
	Avatar   string `json:"avatar"`
}

type Games struct {
	Csgo Csgo `json:"csgo"`
}

type Csgo struct {
	FaceitElo  int    `json:"faceit_elo"`
	LevelLabel string `json:"skill_level_label"`
}

const apiURL = "https://open.faceit.com/data/v4/"
const oldApiURL = "https://api.faceit.com/stats/api/v1/stats/time/users/"

func main() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/impressum/", http.StripPrefix("/impressum", http.FileServer(http.Dir("./impressum"))))
	http.Handle("/datenschutz/", http.StripPrefix("/datenschutz", http.FileServer(http.Dir("./datenschutz"))))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("problem listening to 8080 %v", err)
	}

}
