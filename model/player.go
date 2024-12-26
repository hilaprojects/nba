package model

type Player struct {
	Id            int
	Name          string
	CurrentTeamID int
	Games         []PlayerGameStats
}

type Season struct {
	Id   int
	Year string
}
type Game struct {
	Id       int
	TeamAID  int
	TeamBID  int
	SeasonID int
}
type Team struct {
	Id   int
	Name string
}

type PlayerGameStats struct {
	PlayerID      int
	PlayerName    string
	GameID        int
	Points        int
	Assists       int
	Rebounds      int
	Steals        int
	Blocks        int
	Turnovers     int
	Fouls         int
	MinutesPlayed float32
}
type PlayerSeasonAverage struct {
	PlayerID             int
	PlayerName           string
	Season               int
	PointsPerGame        float32
	AssistsPerGame       float32
	ReboundsPerGame      float32
	StealsPerGame        float32
	BlocksPerGame        float32
	TurnoversPerGame     float32
	FoulsPerGame         float32
	MinutesPlayedPerGame float32
}

type TeamSeasoAverage struct {
	TeamName             string
	TeamID               int
	Season               int
	PointsPerGame        float32
	AssistsPerGame       float32
	ReboundsPerGame      float32
	StealsPerGame        float32
	BlocksPerGame        float32
	TurnoversPerGame     float32
	FoulsPerGame         float32
	MinutesPlayedPerGame float32
}

type GetPlayerGameStatsRequest struct {
	PageSize   int
	PageNumber int
	SeasonYear int
	PlayerID   int
}

type GetPlayerGameStatsResponse struct {
	PlayerStats []PlayerSeasonAverage
}

type GetTeamGameStatsRequest struct {
	PageSize   int
	PageNumber int
	TeamID     int
	SeasonYear int
}

type GetTeamsGameStatsResponse struct {
	TeamStats []TeamSeasoAverage
}

type LogPlayerGameRequest struct {
	PlayerId      int
	GameId        int
	Points        int
	Assists       int
	Rebounds      int
	Steals        int
	Blocks        int
	Turnovers     int
	Fouls         int
	MinutesPlayed float32
}
