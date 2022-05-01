package core

type Game struct {
	Dealer  Cards    `bson:"dealer"`
	Players []Player `bson:"players"`
}

type Player struct {
	Username string `bson:"username"`
	Cards    Cards  `bson:"cards"`
	Stop     bool   `bson:"stop"`
}
