package main

// Player defines the properties of a Player
type Player struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Cards  *CardStack `json:"cards"`
	Class  string     `json:"class"`
	active bool
}

func newPlayer(name string) *Player {
	return &Player{
		ID:    name,
		Name:  name,
		Cards: &CardStack{},
	}
}

// PlayerState is used to send the game state to
// player.
type PlayerState struct {
	Player   Player `json:"player"`
	Opponent Player `json:"opponent"`
	HeapHead Card   `json:"heap_head"`
}

// playerState generates the hand for each player.
func (s *server) playerState(id string) (PlayerState, bool) {
	ps := PlayerState{}
	found := false
	s.game.state()
	for _, p := range s.game.Players {
		switch p.ID {
		case id:
			ps.Player = *p
			if p.active {
				ps.Player.Class = "active"
			}
			found = true
		default:
			opponent := Player{
				ID:   "",
				Name: p.Name,
				Cards: &CardStack{
					// make creates the slice with empty cards
					// because the player should not see or get
					// any information about that cards
					Cards: make([]Card, len(p.Cards.Cards)),
				},
			}
			if p.active {
				opponent.Class = "active"
			}
			for i := range p.Cards.Cards {
				opponent.Cards.Cards[i].Color = "back"
			}
			ps.Opponent = opponent
		}
	}
	ps.HeapHead = s.game.HeapHead
	return ps, found
}
