package game

import (
	"log"
	"math/rand"
)

type CardStruct struct {
	CardType string `json:"cardType"`
	action   func(*GameStruct)
}

// var cardNames = []string{"cat", "exploding", "shuffle", "defuse"}

func GenerateCards() *[]CardStruct {

	catCard := CardStruct{CardType: "cat", action: catCardAction}
	explodingCard := CardStruct{CardType: "exploding", action: explodingCardAction}
	shuffleCard := CardStruct{CardType: "shuffle", action: shuffleCardAction}
	defuseCard := CardStruct{CardType: "defuse", action: defusedCardAction}

	cardsArr := []CardStruct{catCard, defuseCard, shuffleCard, explodingCard}
	// indices := []int{0, 1, 2, 3}

	finalArr := []CardStruct{}

	for i := 0; i < 4; i++ {
		j := rand.Intn(4)
		finalArr = append(finalArr, cardsArr[j])
		// indices = append(indices[:j], indices[j+1:]...)
	}

	finalArr = append(finalArr, cardsArr[rand.Intn(4)])

	return &finalArr
}

func HydrateActionsOnCards(cards *[]CardStruct) {
	log.Println("\n\n\nHyderating cards")

	for i, _ := range *cards {
		log.Println("Card Type: ", (*cards)[i].CardType)
		switch (*cards)[i].CardType {
		case "cat":
			log.Println("Assinging Cat Card Action")
			(*cards)[i].action = catCardAction
		case "shuffle":
			log.Println("Assinging shuffle Card Action")
			(*cards)[i].action = shuffleCardAction
		case "defuse":
			log.Println("Assinging defuse Card Action")
			(*cards)[i].action = defusedCardAction
		case "exploding":
			log.Println("Assinging exploding Card Action")
			(*cards)[i].action = explodingCardAction
		default:
			log.Println("Default executed", (*cards)[i].CardType)
		}
	}
}

func catCardAction(g *GameStruct) {
	log.Println("cat Card")
	removeCard(g)
}

func removeCard(g *GameStruct) {
	length := len(*g.Cards)
	if length == 1 {
		g.Status = "Winner"
	}
	*g.Cards = (*g.Cards)[:length-1]
}

func explodingCardAction(g *GameStruct) {
	log.Println("exploding Card")
	if g.Diffusers > 0 {
		removeCard(g)
		g.Diffusers--
		return
	}
	g.Status = "lost"
}
func shuffleCardAction(g *GameStruct) {
	log.Println("shuffle Card")
	g.Cards = GenerateCards()
	g.Diffusers = 0
}
func defusedCardAction(g *GameStruct) {
	log.Println("diffused Card")
	g.Diffusers++
	removeCard(g)
}
