package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type finalData struct {
	Matches   []OldApiMatch
	PlayerOne PlayerInfo
	PlayerTwo PlayerInfo
}

type OldApiMatch struct {
	MatchID   string `json:"matchId"`
	CreatedAt int    `json:"created_at"`
	//Elo       string `json:"elo"`
}

func findCommonMatches(playerOne, playerTwo string) ([]OldApiMatch, error) {
	playerOneHistory, err := fetchMatchHistory(playerOne, 0)
	if err != nil {
		return nil, err
	}
	playerTwoHistory, err := fetchMatchHistory(playerTwo, 0)
	if err != nil {
		return nil, err
	}
	if len(playerOneHistory) > len(playerTwoHistory) {
		return mapCommonMatches(playerTwoHistory, playerOneHistory)
	}
	return mapCommonMatches(playerOneHistory, playerTwoHistory)

}

func mapCommonMatches(playerOne, playerTwo []OldApiMatch) ([]OldApiMatch, error) {
	matches := map[string]bool{}
	for _, match := range playerOne {
		matches[match.MatchID] = true
	}
	var commonMatches []OldApiMatch
	for _, match := range playerTwo {
		_, exists := matches[match.MatchID]
		if exists {
			commonMatches = append(commonMatches, match)
		}
	}
	return commonMatches, nil
}

func fetchMatchHistory(playerID string, page int) ([]OldApiMatch, error) {
	req, err := http.NewRequest("GET", oldApiURL+playerID+"/games/csgo?size=2000&page="+strconv.Itoa(page), nil)
	if err != nil {
		return []OldApiMatch{}, err
	}
	req.Header.Set("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []OldApiMatch{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		return []OldApiMatch{}, fmt.Errorf("Got HTTP status %d, want 200 < status < 299", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var jsonData []OldApiMatch
	if err := decoder.Decode(&jsonData); err != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
		return []OldApiMatch{}, err
	}

	if len(jsonData) > 1999 {
		tail, err := fetchMatchHistory(playerID, page+1)
		if err != nil {
			return []OldApiMatch{}, err
		}
		jsonData = append(jsonData, tail...)
	}
	return jsonData, nil
}

func fetchPlayerData(nickname string) (Player, error) {
	req, err := http.NewRequest("GET", apiURL+"players?nickname="+nickname, nil)
	if err != nil {
		return Player{}, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FACEIT_TOKEN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Player{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		return Player{}, fmt.Errorf("Got HTTP status %d, want 200 < status < 299", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var data Player
	if err := decoder.Decode(&data); err != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
		return Player{}, err
	}
	return data, nil
}
