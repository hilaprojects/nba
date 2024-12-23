package postgres

import (
	"nba/model"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	LogPlayerGame(game model.PlayerGameStat) error
	GetTeamPlayers(teamID int) ([]model.Player, error)
	GetPlayer(playerId int) (model.Player, error)
	GetPlayers(pageNumber int, pageSize int) ([]model.Player, int64, error)
	GetPlayerGames(playerId int) ([]model.PlayerGameStat, error)
	SavePlayer(player model.Player) error
	SaveTeam(team model.Team) error
	GetTeams() ([]model.Team, error)
}

type PlayerRepositoryStruct struct {
	db *gorm.DB
}

// GetPlayer implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetPlayer(playerId int) (model.Player, error) {
	panic("unimplemented")
}

// GetPlayerGames implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetPlayerGames(playerId int) ([]model.PlayerGameStat, error) {
	panic("unimplemented")
}

// GetTeamPlayers implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetTeamPlayers(teamID int) ([]model.Player, error) {
	var players []model.Player
	err := p.db.Debug().Where("team_id = ?", teamID).Find(&players).Error
	if err != nil {
		return nil, err
	}
	return players, nil
}

// GetPlayers implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetPlayers(pageNumber int, pageSize int) ([]model.Player, int64, error) {
	var count int64
	err := p.db.Model(&model.Player{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (pageNumber - 1) * pageSize
	var players []model.Player

	err = p.db.Debug().
		Offset(offset).
		Limit(pageSize).
		Find(&players).Error

	if err != nil {
		return nil, 0, err
	}
	return players, count, nil
}

// LogPlayerGame implements PlayerRepository.
func (p *PlayerRepositoryStruct) LogPlayerGame(game model.PlayerGameStat) error {
	err := p.db.Debug().Save(&game).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PlayerRepositoryStruct) SavePlayer(player model.Player) error {
	err := p.db.Debug().Save(&player).Error
	if err != nil {
		return err
	}
	return nil
}
func (p *PlayerRepositoryStruct) SaveTeam(team model.Team) error {
	err := p.db.Debug().Save(&team).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PlayerRepositoryStruct) GetTeams() ([]model.Team, error) {
	var teams []model.Team
	err := p.db.Debug().Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &PlayerRepositoryStruct{
		db: db,
	}
}
