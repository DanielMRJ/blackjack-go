package main

import (
	"fmt"
	"math/rand"
	"time"
)

var suits = []string{"♠", "♥", "♦", "♣"}
var values = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

type Card struct {
	Suit  string
	Value string
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Value, c.Suit)
}

type Deck []Card

func NewDeck() Deck {
	deck := Deck{}
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

func (d Deck) Shuffle() Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
	return d
}

func (d *Deck) Draw() Card {
	card := (*d)[0]
	*d = (*d)[1:]
	return card
}
