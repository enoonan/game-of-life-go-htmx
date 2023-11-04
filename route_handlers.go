package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func next(c echo.Context) error {
	gameId := c.Param("id")
	gameIdInt, _ := strconv.Atoi(gameId)
	game := FetchLatestGameState(gameIdInt)
	newGame := runCycle(game)

	if newGame.Dead {
		KillGame(gameIdInt)
	}
	newGame.Name = game.Name
	SaveGameState(gameIdInt, newGame.Board)
	data := pageData()
	data.Game = newGame
	return c.Render(http.StatusOK, "game-board", data)
}

func home(c echo.Context) error {
	game := FetchLatestGame()
	data := pageData()
	data.Game = game
	return c.Render(http.StatusOK, "index", data)
}

func show(c echo.Context) error {
	gameId := c.Param("id")
	gameIdInt, _ := strconv.Atoi(gameId)
	game := FetchSingleGame(gameIdInt)
	state := FetchLatestGameState(gameIdInt)
	game.Board = state.Board
	data := pageData()
	data.Game = game
	return c.Render(http.StatusOK, "index", data)
}

func start(c echo.Context) error {
	game := MakeNewGame()
	pageData := pageData()
	pageData.Game = game

	c.Response().Header().Set("HX-Location", host+"/games/"+strconv.Itoa(game.Id))
	return c.Render(http.StatusOK, "index", pageData)
}

func listGames(c echo.Context) error {
	games := FetchGames()
	data := pageData()
	data.Games = games
	return c.Render(http.StatusOK, "games-list", data)
}

func getRunner(c echo.Context) error {
	gameId := c.Param("id")
	gameIdInt, _ := strconv.Atoi(gameId)
	game := Game{Id: gameIdInt}
	data := pageData()
	data.Game = game
	return c.Render(http.StatusOK, "game-controls-running", data)
}

func cancelRunner(c echo.Context) error {
	gameId, _ := strconv.Atoi(c.Param("id"))
	data := pageData()
	data.Game = Game{Id: gameId}
	return c.Render(http.StatusOK, "game-controls-default", data)
}
