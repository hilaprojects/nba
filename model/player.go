package model

type Player struct {
	Id     int              `gorm:"primaryKey;autoIncrement"` // Primary key with auto-increment
	Name   string           `gorm:"type:varchar(100);not null"`
	TeamID int              `gorm:"not null;index;constraint:OnDelete:CASCADE"`                       // Foreign key for the relationship with Team, with cascade delete
	Games  []PlayerGameStat `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,onDelete:CASCADE"` // One-to-many relationship with cascade delete
}

type PlayerGameStat struct {
	PlayerID      int // Foreign key for the relationship with Player
	GameID        int // Foreign key for the relationship with Game
	Points        int
	Assists       int
	Rebounds      int
	Steals        int
	Blocks        int
	Turnovers     int
	Fouls         int
	MinutesPlayed float32
}
type Game struct {
	Id   int    `gorm:"primaryKey;autoIncrement"` // Primary key with auto-increment
	Date string `gorm:"type:date;not null"`       // Date of the game
	Team []Team `gorm:"many2many:game_teams"`     // Many-to-many relationship with Team
}
type Team struct {
	Id   int    `gorm:"primaryKey;autoIncrement"`          // Primary key with auto-increment
	Name string `gorm:"type:varchar(100);unique;not null"` // Unique team name
}
type PlayerSeasonAverage struct {
	PlayerID             int
	TeamID               int
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
	PointsPerGame        float32
	AssistsPerGame       float32
	ReboundsPerGame      float32
	StealsPerGame        float32
	BlocksPerGame        float32
	TurnoversPerGame     float32
	FoulsPerGame         float32
	MinutesPlayedPerGame float32
}

type GetPlayersGameStatsRequest struct {
	PageSize   int
	PageNumber int
}

type GetPlayersGameStatsResponse struct {
	PlayerStats []PlayerSeasonAverage
}

type GetTeamsGameStatsRequest struct {
	PageSize   int
	PageNumber int
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
