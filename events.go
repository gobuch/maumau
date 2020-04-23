package main

import "github.com/google/uuid"

// addCardGameToStack takes a CardStack and add that to the
// game. The card deck for the game need to be a CardStack.
func addCardGameToStack(cs *CardStack) Event {
	return func(g *GameState) {
		g.Stack.Cards = append(cs.Cards, g.Stack.Cards...)
	}
}

// addPlayer takes a Player and adds him/her to the game
func addPlayer(p *Player) Event {
	return func(g *GameState) {
		p.Cards = &CardStack{}
		if p.ID == "" {
			p.ID = uuid.New().String()
		}
		g.Players = append(g.Players, p)
	}
}

// setNextPlayer changes the active player
func setNextPlayer(p *Player) Event {
	return func(g *GameState) {
		next, _ := g.nextPlayer(p.ID)
		g.setActivePlayer(next.ID)
	}
}

// serveGame serves the cards from the stack to the players
func serveGame() Event {
	return func(g *GameState) {
		// a new emtpy hand for every player
		for _, p := range g.Players {
			p.Cards = &CardStack{Cards: []Card{}}
			p.active = true
		}
		for i := 1; i <= g.NrCards; i++ {
			for j := range g.Players {
				g.Players[j].Cards.push(g.Stack.pop())
			}
		}
		g.Heap.push(g.Stack.pop())
	}
}

// takeCardFromStack lets a player take a card
func takeCardFromStack(p *Player) Event {
	return func(g *GameState) {
		p.Cards.push(g.Stack.pop())
	}
}

// playCardToHeap removes the card from the hand of a player
// and adds the card to the heap
func playCardToHeap(p *Player, i int) Event {
	return func(g *GameState) {
		g.Heap.push(p.Cards.take(i))
		head := g.Heap.peek()
		next, _ := g.nextPlayer(p.ID)
		for i := 0; i < head.SkipPlayers; i++ {
			next, _ = g.nextPlayer(next.ID)
		}
		g.setActivePlayer(next.ID)
	}
}

// removeCardsFromHeap removes all cards from the heap. Just the
// top card stays.
// This event is used, when there are no more cards on the stack.
// In that case all played cards of the heap are removed from the
// heap and added to the stack.
// The stack is shuffeld and added to the card game.
func removeCardsFromHeap() Event {
	return func(g *GameState) {
		g.Heap.Cards = []Card{g.Heap.peek()}
	}
}
