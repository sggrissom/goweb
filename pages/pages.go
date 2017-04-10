package main

import (
	"fmt"
)

type page struct {
	Path     string
	Template string
	Title    string
	Display  string
	JS       []string
	InNav    bool
	DataFn   func() map[string]interface{}
}

type player struct {
	First       string
	Last        string
	GamesPlayed int
}

func RegisterPages() {
	pages = []page{
		page{
			Path:     "players",
			Title:    "playas",
			Display:  "Players",
			JS:       []string{"players.js"},
			InNav:    true,
			DataFn:   PlayersPage,
		},
		page{Path: "home", Display: "Homey"},
		page{Path: "settings", Display: "Settings", InNav: true},
	}
}

func PlayersPage() map[string]interface{} {
	fmt.Printf("show the playas\n")
	m := make(map[string]interface{})

	m["players"] = []player{
		player{First: "Aria", Last: "Grissom", GamesPlayed: 4},
		player{First: "Steven", Last: "Grissom", GamesPlayed: 74},
		player{First: "a", Last: "person", GamesPlayed: 1},
		player{First: "that", Last: "guy", GamesPlayed: 0},
	}
	return m
}
