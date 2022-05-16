package core

const (
	GameDealerField         = "dealer"
	GamePlayersField        = "players"
	GamePlayerUsernameField = "username"
	GamePlayerCardsField    = "cards"
	GamePlayerStopField     = "stop"
	GamePlayerBustedField   = "busted"
)

type Game struct {
	Dealer  Cards    `bson:"dealer"`
	Players []Player `bson:"players"`
}

type Player struct {
	Username string `bson:"username"`
	Cards    Cards  `bson:"cards"`
	Stop     bool   `bson:"stop"`
	Busted   bool   `bson:"busted"`
}
