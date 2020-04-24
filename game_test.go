package main

import (
	"reflect"
	"testing"
)

func TestGameEvents(t *testing.T) {
	g := newGame()
	cardGame := CardGame()
	g.event(addCardGameToStack(cardGame))
	g.event(addPlayer(newPlayer("Max")))
	g.event(addPlayer(newPlayer("Maja")))
	g.event(serveGame())
	g.state() // State can be called more then once
	g.state()
	// check if there are cards in the deck
	if g.Stack.len() == 0 {
		t.Error("stack is empty some cards should be added to the game")
	}
	// check if the first player is Max
	if g.Players[0].Name != "Max" {
		t.Errorf("first player should be Max: got %#v", g.Players[0])
	}
	// check if there is a ID for the first player
	if g.Players[0].ID == "" {
		t.Error("Player 0 has no ID")
	}
	// check if the first player has card
	if g.Players[0].Cards.len() != g.NrCards {
		t.Errorf("Player 0 had %d cards! Expect: %d", g.Players[0].Cards.len(), g.NrCards)
		t.Errorf("%#v", g.Players[0].Cards)
	}

	// remember the top card
	topCard := g.Stack.peek()
	// The first player should take a card from the stack
	g.event(takeCardFromStack(g.Players[0]))
	g.state()
	pTopCard := g.Players[0].Cards.peek()
	// check if the player took the topCard
	if !reflect.DeepEqual(topCard, pTopCard) {
		t.Error("popCardFromStack does not serve the right card")
	}
	// player plays one card to the heap
	pFirstCard := g.Players[0].Cards.Cards[0]
	g.event(playCardToHeap(g.Players[0], 0))
	g.state()

	hTopCard := g.Heap.peek()
	// the Card need an ID
	if pFirstCard.ID == "" {
		t.Error("pFirstCard has no ID")
	}
	// check if the right card is played
	if pFirstCard.ID != hTopCard.ID {
		t.Error("pushCardToHeap does not push the right card to the heap")
	}
}
