package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

func db() *sql.DB {
	dbConn, err := sql.Open("sqlite3", "./game_of_life.db")
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func InitializeDB() {
	conn := db()
	gameSql := `CREATE TABLE IF NOT EXISTS games(
    id INTEGER PRIMARY KEY,
    name STRING,
		dimension INTEGER,
		dead INTEGER DEFAULT 0
);`
	_, err := conn.Exec(gameSql)
	if err != nil {
		log.Fatal(err)
	}

	gameStateSql := `CREATE TABLE IF NOT EXISTS game_states(
		id INTEGER PRIMARY KEY,
		game_id INTEGER,
		state TEXT
	);`
	_, err = conn.Exec(gameStateSql)
	if err != nil {
		log.Fatal(err)
	}
}

func FetchLatestGame() Game {
	conn := db()
	defer conn.Close()

	var id int
	var name string
	err := conn.QueryRow("SELECT id, name FROM games ORDER BY id DESC LIMIT 1").Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return MakeNewGame()
		}

		log.Fatal(err)
	}

	state := FetchLatestGameState(id)

	game := Game{Id: id, Board: state.Board, Name: name, Dimension: len(state.Board)}
	return game
}

func MakeNewGame() Game {
	game := Game{Board: initBoard(dimension), Dimension: dimension}
	conn := db()
	defer conn.Close()
	var gameCount int
	err := conn.QueryRow("SELECT COUNT(*) FROM games").Scan(&gameCount)
	if err != nil {
		log.Fatal(err)
	}
	gameCount++
	gameName := fmt.Sprintf("Game of Life #%d", gameCount)
	result, err2 := conn.Exec("INSERT INTO games (name, dimension) VALUES (?, ?)", gameName, game.Dimension)
	if err2 != nil {
		log.Fatal(err2)
	}
	id, _ := result.LastInsertId()
	game.Id = int(id)
	game.Dead = false
	game.Name = gameName
	game.Dimension = dimension

	SaveGameState(game.Id, game.Board)
	return game
}

func FetchLatestGameState(gameId int) Game {
	conn := db()
	defer conn.Close()

	var state string
	err := conn.QueryRow("SELECT state FROM game_states WHERE game_id = ? ORDER BY id DESC LIMIT 1", gameId).Scan(&state)
	if err != nil {
		log.Fatal(err)
	}
	var matrix [][]int
	err = json.Unmarshal([]byte(state), &matrix)
	if err != nil {
		log.Fatal(err)
	}

	return Game{Id: gameId, Board: matrix, Dimension: len(matrix)}
}

func SaveGameState(gameId int, state [][]int) {
	conn := db()
	defer conn.Close()
	stateStr, _ := json.Marshal(state)
	_, err := conn.Exec("INSERT INTO game_states (game_id, state) VALUES (?, ?)", gameId, stateStr)
	if err != nil {
		log.Fatal(err)
	}
}

func KillGame(gameId int) {
	conn := db()
	defer conn.Close()
	_, err := conn.Exec("UPDATE games SET dead = 1 WHERE id = ?", gameId)
	if err != nil {
		log.Fatal(err)
	}
}

func FetchSingleGame(id int) Game {
	conn := db()
	defer conn.Close()
	var name string
	var dimension int
	var dead int

	err := conn.QueryRow("SELECT name, dimension, dead FROM games WHERE id = ?", id).Scan(&name, &dimension, &dead)
	if err != nil {
		log.Fatal(err)
	}

	game := Game{Id: id, Name: name, Dimension: dimension, Dead: dead == 1}
	return game
}

func FetchGames() []Game {
	conn := db()
	defer conn.Close()

	rows, err := conn.Query("SELECT id, name, dimension, dead FROM games ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	games := []Game{}
	for rows.Next() {
		var id int
		var name string
		var dimension int
		var dead int
		err = rows.Scan(&id, &name, &dimension, &dead)
		if err != nil {
			log.Fatal(err)
		}
		games = append(games, Game{Id: id, Name: name, Dimension: dimension, Dead: dead == 1})
	}
	return games
}
