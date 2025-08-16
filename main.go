package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

type Game struct {
	Player   []int
	Dealer   []int
	Message  string
	GameOver bool
}

// cartas posibles (simplificado: As vale 11, J/Q/K = 10)
var deck = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 11}

var game Game

func drawCard() int {
	rand.Seed(time.Now().UnixNano())
	return deck[rand.Intn(len(deck))]
}

func newGame() {
	game = Game{
		Player:   []int{drawCard(), drawCard()},
		Dealer:   []int{drawCard()},
		Message:  "Tu turno, elige una acción",
		GameOver: false,
	}
}

func total(hand []int) int {
	sum := 0
	for _, c := range hand {
		sum += c
	}
	return sum
}

func handler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")

	if action == "new" || game.Player == nil {
		newGame()
	} else if !game.GameOver {
		switch action {
		case "hit":
			game.Player = append(game.Player, drawCard())
			if total(game.Player) > 21 {
				game.Message = "Te pasaste de 21 😢"
				game.GameOver = true
			}
		case "stand":
			// turno del dealer
			for total(game.Dealer) < 17 {
				game.Dealer = append(game.Dealer, drawCard())
			}
			// decidir ganador
			playerScore := total(game.Player)
			dealerScore := total(game.Dealer)

			if dealerScore > 21 || playerScore > dealerScore {
				game.Message = "¡Ganaste! 🎉"
			} else if playerScore == dealerScore {
				game.Message = "Empate 🤝"
			} else {
				game.Message = "Perdiste 😭"
			}
			game.GameOver = true
		}
	}

	tmpl.Execute(w, game)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
