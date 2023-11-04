package main

import "html/template"

type PageData struct {
	Title string
	Game  Game
	Host  string
	Games []Game
}

type Template struct {
	templates *template.Template
}

type Game struct {
	Id        int
	Board     [][]int
	Dimension int
	Dead      bool
	Name      string
}
