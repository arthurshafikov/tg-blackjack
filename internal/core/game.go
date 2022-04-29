package core

type Game struct {
	Deck         *Deck        `bson:"deck"`
	DealerHand   Cards        `bson:"dealer_hand"`
	PlayersHands []PlayerHand `bson:"players_hands"`
}

type PlayerHand struct {
	Username string `bson:"username"`
	Cards    Cards  `bson:"cards"`
	Stop     bool   `bson:"stop"`
}
