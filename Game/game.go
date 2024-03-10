package game

// Core Functionality:
// 1. User should login
// 2. User should start the game
//	2.1 User should be presented with a new Game Object with 5 card.
// 3. User should draw cards one-by-one and face the subjected actions for that card
// 4. If the user is manages to draw all 5 cards he/she wins the game and a point is added to his/her account

type GameStruct struct {
	Cards     *[]CardStruct `json:"cards"`
	Diffusers int           `json:"diffusers"`
	Status    string        `json:"status"`
}

func StartNewGame() *GameStruct {
	cards := GenerateCards()
	return &GameStruct{Cards: cards, Status: "playing"}
}

func (G *GameStruct) DrawCard() {
	i := len(*G.Cards)
	if i > 0 {
		(*G.Cards)[i-1].action(G)
	}
}
