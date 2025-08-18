package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

type Game struct {
	PlayerCards []string `json:"playerCards"`
	DealerCards []string `json:"dealerCards"`
	PlayerScore int      `json:"playerScore"`
	DealerScore int      `json:"dealerScore"`
	Message     string   `json:"message"`
}

var deck []string
var player []string
var dealer []string

// Crear y barajar mazo
func newDeck() []string {
	suits := []string{"♠", "♥", "♦", "♣"}
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	var d []string
	for _, s := range suits {
		for _, v := range values {
			d = append(d, v+s)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
	return d
}

// Calcular puntaje
func score(cards []string) int {
	total := 0
	aces := 0
	for _, c := range cards {
		val := c[:len(c)-3] // número/letra (quitando símbolo)
		switch val {
		case "A":
			total += 11
			aces++
		case "K", "Q", "J", "10":
			total += 10
		default:
			n, _ := strconv.Atoi(val)
			total += n
		}
	}
	// Ajustar ases
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// Iniciar partida
func startHandler(w http.ResponseWriter, r *http.Request) {
	deck = newDeck()
	player = []string{deck[0], deck[2]}
	dealer = []string{deck[1], deck[3]}
	deck = deck[4:]

	game := Game{
		PlayerCards: player,
		DealerCards: []string{dealer[0], "??"}, // ocultar una carta dealer
		PlayerScore: score(player),
		DealerScore: 0,
		Message:     "Juego iniciado 🚀",
	}
	json.NewEncoder(w).Encode(game)
}

// Pedir carta
func hitHandler(w http.ResponseWriter, r *http.Request) {
	if len(deck) == 0 {
		deck = newDeck()
	}
	player = append(player, deck[0])
	deck = deck[1:]

	game := Game{
		PlayerCards: player,
		DealerCards: []string{dealer[0], "??"},
		PlayerScore: score(player),
		Message:     "Pediste carta",
	}

	if game.PlayerScore > 21 {
		game.Message = "Te pasaste 💀 ¡Perdiste!"
		game.DealerCards = dealer
		game.DealerScore = score(dealer)
	}
	json.NewEncoder(w).Encode(game)
}

// Plantarse
func standHandler(w http.ResponseWriter, r *http.Request) {
	for score(dealer) < 17 {
		dealer = append(dealer, deck[0])
		deck = deck[1:]
	}

	playerScore := score(player)
	dealerScore := score(dealer)

	msg := ""
	if dealerScore > 21 || playerScore > dealerScore {
		msg = "¡Ganaste! 🎉"
	} else if playerScore < dealerScore {
		msg = "Perdiste 😢"
	} else {
		msg = "Empate 🤝"
	}

	game := Game{
		PlayerCards: player,
		DealerCards: dealer,
		PlayerScore: playerScore,
		DealerScore: dealerScore,
		Message:     msg,
	}
	json.NewEncoder(w).Encode(game)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/hit", hitHandler)
	http.HandleFunc("/stand", standHandler)

	log.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
