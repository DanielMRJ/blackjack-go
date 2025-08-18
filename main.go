package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Card struct {
	Value int    // 2-10, 10 para J/Q/K, 11 para As
	Face  string // "A","2"... "K"
	Suit  string // ♠ ♥ ♦ ♣
}

type GameState struct {
	PlayerCards []string `json:"playerCards"`
	DealerCards []string `json:"dealerCards"`
	PlayerScore int      `json:"playerScore"`
	DealerScore int      `json:"dealerScore"`
	Message     string   `json:"message"`
	GameOver    bool     `json:"gameOver"`
}

var (
	tpl    = template.Must(template.ParseFiles("templates/index.html"))
	deck   []Card
	player []Card
	dealer []Card
	ended  bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/hit", hitHandler)
	http.HandleFunc("/stand", standHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// ------------ Handlers ------------

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.Execute(w, nil)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	newGame()
	// Reparte: jugador 2, dealer 2 (mostramos 1 y ocultamos 1 hasta 'stand')
	player = append(player, draw(), draw())
	dealer = append(dealer, draw(), draw())
	ended = false
	writeJSON(w, state(false), http.StatusOK)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	if ended || len(player) == 0 {
		writeJSON(w, GameState{Message: "Primero pulsa Iniciar"}, http.StatusOK)
		return
	}
	player = append(player, draw())
	ps := score(player)
	if ps > 21 {
		ended = true
		st := state(true)
		st.Message = "¡Te pasaste! Pierdes 💀"
		writeJSON(w, st, http.StatusOK)
		return
	}
	writeJSON(w, state(false), http.StatusOK)
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	if len(player) == 0 {
		writeJSON(w, GameState{Message: "Primero pulsa Iniciar"}, http.StatusOK)
		return
	}
	// Dealer roba hasta 17 o más
	for score(dealer) < 17 {
		dealer = append(dealer, draw())
	}
	ended = true
	ps, ds := score(player), score(dealer)
	msg := ""
	switch {
	case ds > 21 || ps > ds:
		msg = "¡Ganaste! 🎉"
	case ps < ds:
		msg = "La banca gana 😕"
	default:
		msg = "Empate 🤝"
	}
	st := state(true)
	st.Message = msg
	writeJSON(w, st, http.StatusOK)
}

// ------------ Lógica de juego ------------

func newGame() {
	deck = makeDeck()
	shuffle(deck)
	player = []Card{}
	dealer = []Card{}
}

func makeDeck() []Card {
	suits := []string{"♠", "♥", "♦", "♣"}
	faces := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	d := make([]Card, 0, 52)
	for _, s := range suits {
		for _, f := range faces {
			val := 0
			switch f {
			case "A":
				val = 11
			case "J", "Q", "K":
				val = 10
			default:
				n, _ := strconv.Atoi(f)
				val = n
			}
			d = append(d, Card{Value: val, Face: f, Suit: s})
		}
	}
	return d
}

func shuffle(d []Card) {
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
}

func draw() Card {
	if len(deck) == 0 {
		deck = makeDeck()
		shuffle(deck)
	}
	c := deck[0]
	deck = deck[1:]
	return c
}

func score(hand []Card) int {
	total, aces := 0, 0
	for _, c := range hand {
		total += c.Value
		if c.Face == "A" {
			aces++
		}
	}
	for total > 21 && aces > 0 {
		total -= 10 // As de 11 -> 1
		aces--
	}
	return total
}

func renderCards(hand []Card) []string {
	out := make([]string, len(hand))
	for i, c := range hand {
		out[i] = c.Face + c.Suit
	}
	return out
}

func state(revealDealer bool) GameState {
	ds := 0
	var dc []string
	if revealDealer {
		ds = score(dealer)
		dc = renderCards(dealer)
	} else {
		// muestra primera y oculta la segunda
		if len(dealer) == 0 {
			dc = []string{}
		} else if len(dealer) == 1 {
			dc = []string{dealer[0].Face + dealer[0].Suit}
		} else {
			dc = []string{dealer[0].Face + dealer[0].Suit, "??"}
		}
	}
	return GameState{
		PlayerCards: renderCards(player),
		DealerCards: dc,
		PlayerScore: score(player),
		DealerScore: ds,
		Message:     "",
		GameOver:    ended,
	}
}

// ------------ Util ------------

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
