package core

type Game struct {
	DealerHand Cards    `bson:"dealer_hand"`
	Players    []Player `bson:"players"`
}

type Player struct {
	Username string `bson:"username"`
	Cards    Cards  `bson:"cards"`
	Stop     bool   `bson:"stop"`
}
