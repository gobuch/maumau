package main

// event is a function, which takes a pointer to the game
// When the game calculates the state all events are called
type event func(g *GameState)

// Game defines the logic of the maumau game
type Game struct {
	GameState
	events     []event
	redoEvents []event
}

// GameState includes all properties of the game, which are having
// a state. The is just changed by an event.
type GameState struct {
	Stack        *CardStack `json:"stack"`
	Heap         *CardStack `json:"heap"`
	HeapHead     Card       `json:"heap_head"`
	Players      []*Player  `json:"players"`
	ActivePlayer int        `json:"active_player"`
	NrCards      int        `json:"nr_cards"`
}

func newGame() *Game {
	return &Game{
		GameState: GameState{
			Stack:        &CardStack{},
			Heap:         &CardStack{},
			ActivePlayer: 0,
			NrCards:      6,
		},
	}
}

// event adds the event e to the game
func (g *Game) event(e event) {
	// every new event clears the redo slice
	g.redoEvents = []event{}
	g.events = append(g.events, e)
}

// init clears the internal state
func (g *Game) init() {
	g.Stack = &CardStack{}
	g.Heap = &CardStack{}
	g.Players = []*Player{}
}

// state executes all events and creates a new state for
// the Stack, Heap and the Players
func (g *Game) state() {
	g.init()
	for _, e := range g.events {
		e(&g.GameState)
	}
	g.HeapHead = g.Heap.peek()
}

// player returns the player matching to the given ID.
func (g *GameState) player(id string) (*Player, bool) {
	for _, p := range g.Players {
		if id == p.ID {
			return p, true
		}
	}
	// if there is no player with id
	return nil, false
}

// nextPlayer takes the current player ID and returns the
// next player at the table.
func (g *GameState) nextPlayer(id string) (*Player, bool) {
	found := -1
	// index of the current player
	for i, p := range g.Players {
		if id == p.ID {
			found = i
			break
		}
	}
	// given id does not exist
	if found == -1 {
		return nil, false
	}
	next := found + 1
	if len(g.Players)-1 == found {
		// if current player is the last in the slice
		// the first player is next
		next = 0
	}
	return g.Players[next], true
}

// SetActivePlayer sets the active player
// all other players will set es not active.
func (g *GameState) setActivePlayer(id string) bool {
	// first check if there is a player with the id
	_, ok := g.player(id)
	if !ok {
		return false
	}
	for _, p := range g.Players {
		if id == p.ID {
			p.active = true
		} else {
			p.active = false
		}
	}
	return true
}
