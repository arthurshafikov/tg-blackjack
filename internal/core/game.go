package core

type Game struct {
	Deck         Cards                 `bson:"deck"`
	DealerHand   Cards                 `bson:"dealer_hand"`
	PlayersHands map[string]PlayerHand `bson:"players_hand"`
	IsFinished   bool                  `bson:"is_finished"`
}

type PlayerHand struct {
	Cards Cards `bson:"cards"`
	Stop  bool  `bson:"stop"`
}
