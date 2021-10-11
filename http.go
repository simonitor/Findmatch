package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var commonMatches []OldApiMatch
	var playerOne, playerTwo PlayerInfo

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("parse froms: %v", err), http.StatusInternalServerError)
			return
		}
		nicknamePlayerOne := r.PostFormValue("playerOne")
		nicknamePlayerTwo := r.PostFormValue("playerTwo")

		log.Printf("0 = %s 1 = %s", nicknamePlayerOne, nicknamePlayerTwo)
		pi1, err := fetchPlayerData(nicknamePlayerOne)
		if err != nil {
			http.Error(w, fmt.Sprintf("fetch data: %v", err), http.StatusInternalServerError)
			return
		}
		pi2, err := fetchPlayerData(nicknamePlayerTwo)
		if err != nil {
			http.Error(w, fmt.Sprintf("fetch data: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("%v %v", pi1.PlayerID, pi2.PlayerID)

		playerOne = PlayerInfo{
			PlayerID:   pi1.PlayerID,
			Elo:        pi1.Games.Csgo.FaceitElo,
			Nickname:   nicknamePlayerOne,
			Avatar:     pi1.Avatar,
			SkillLevel: pi1.Games.Csgo.LevelLabel,
		}
		playerTwo = PlayerInfo{
			PlayerID:   pi2.PlayerID,
			Elo:        pi2.Games.Csgo.FaceitElo,
			Nickname:   nicknamePlayerTwo,
			Avatar:     pi2.Avatar,
			SkillLevel: pi2.Games.Csgo.LevelLabel,
		}
		if playerOne.PlayerID != "" && playerTwo.PlayerID != "" {
			commonMatches, err = findCommonMatches(playerOne.PlayerID, playerTwo.PlayerID)
			if err != nil {
				http.Error(w, fmt.Sprintf("fetch data: %v", err), http.StatusInternalServerError)
				return
			}
		}

		//log.Print(commonMatches)
	}

	log.Println(r.Method, r.RequestURI)
	funcs := template.FuncMap{
		"toDate": func(timestamp int) string {
			if timestamp == 0 {
				return "unknown"
			}
			date := time.Unix(int64(timestamp/1000), 0)
			return date.Format("02.01.2006 - 15:04")
		},
	}

	t, err := template.New("index.html").Funcs(funcs).ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("parse files: %v", err), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, finalData{Matches: commonMatches, PlayerOne: playerOne, PlayerTwo: playerTwo}); err != nil {
		http.Error(w, fmt.Sprintf("execute template: %v", err), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
