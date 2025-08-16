package main

import (
	"fmt"
)

func main() {
	deck := NewDeck().Shuffle()

	player := Player{Name: "Jugador"}
	dealer := Player{Name: "Dealer"}

	// Repartir dos cartas a cada uno
	player.AddCard(deck.Draw())
	player.AddCard(deck.Draw())
	dealer.AddCard(deck.Draw())
	dealer.AddCard(deck.Draw())

	fmt.Println("Tus cartas:", player.ShowHand(false))
	fmt.Println("Cartas del dealer:", dealer.ShowHand(true))

	// Turno del jugador
	var action string
	for player.HandValue() < 21 {
		fmt.Print("¿Pedir carta (h) o plantarse (s)? ")
		fmt.Scanln(&action)
		if action == "h" {
			player.AddCard(deck.Draw())
			fmt.Println("Tus cartas:", player.ShowHand(false))
		} else {
			break
		}
	}

	playerTotal := player.HandValue()
	if playerTotal > 21 {
		fmt.Println("Te pasaste! Puntaje:", playerTotal, "Dealer gana 😢")
		return
	}

	// Turno del dealer
	fmt.Println("\nTurno del dealer...")
	fmt.Println("Cartas del dealer:", dealer.ShowHand(false))
	for dealer.HandValue() < 17 {
		dealer.AddCard(deck.Draw())
		fmt.Println("Cartas del dealer:", dealer.ShowHand(false))
	}

	dealerTotal := dealer.HandValue()
	fmt.Println("\nPuntaje final - Jugador:", playerTotal, "Dealer:", dealerTotal)

	if dealerTotal > 21 || playerTotal > dealerTotal {
		fmt.Println("¡Ganaste! 🎉")
	} else if playerTotal < dealerTotal {
		fmt.Println("Dealer gana 😢")
	} else {
		fmt.Println("Empate 😐")
	}
}
