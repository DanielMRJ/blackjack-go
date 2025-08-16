package main

import "fmt"

type Player struct {
	Name string
	Hand []Card
}

func (p *Player) AddCard(c Card) {
	p.Hand = append(p.Hand, c)
}

func (p Player) HandValue() int {
	total := 0
	aces := 0
	for _, c := range p.Hand {
		switch c.Value {
		case "J", "Q", "K":
			total += 10
		case "A":
			total += 11
			aces++
		default:
			total += atoi(c.Value)
		}
	}
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

func (p Player) ShowHand(hidden bool) string {
	if hidden {
		return fmt.Sprintf("%s ??", p.Hand[0])
	}
	s := ""
	for _, c := range p.Hand {
		s += c.String() + " "
	}
	return s
}

func atoi(str string) int {
	var num int
	fmt.Sscan(str, &num)
	return num
}
