package postgres

import (
	"database/sql"
	"fmt"
	"nba/model"
)

type PlayerRepository interface {
	LogPlayerGame(game model.PlayerGameStats) error
	GetPlayerGamesBySeason(playerID int, season int) ([]model.PlayerGameStats, error)
	GetTeamPlayersBySeason(teamID int, season int) ([]model.PlayerGameStats, error)
	GetPlayer(playerId int) (model.Player, error)
	GetGame(gameId int) (model.Game, error)
	GetTeam(teamId int) (model.Team, error)
}

type PlayerRepositoryStruct struct {
	db *sql.DB
}

// GetTeam implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetTeam(teamId int) (model.Team, error) {
	var team model.Team
	err := p.db.QueryRow("SELECT * from team where id = $1", teamId).Scan(team).Error
	if err != nil {
		return model.Team{}, nil
	}
	return team, nil
}

// GetGame implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetGame(gameId int) (model.Game, error) {
	var game model.Game
	err := p.db.QueryRow("SELECT * from game where id = $1", gameId).Scan(game).Error
	if err != nil {
		return model.Game{}, nil
	}
	return game, nil
}

// GetPlayer implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetPlayer(playerId int) (model.Player, error) {
	var player model.Player
	err := p.db.QueryRow("SELECT * from player where id = $1", playerId).Scan(player).Error
	if err != nil {
		return model.Player{}, nil
	}
	return player, nil
}

// GetPlayerGames implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetPlayerGames(playerId int) ([]model.PlayerGameStats, error) {
	panic("unimplemented")
}

// GetTeamPlayers implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetTeamPlayersBySeason(teamID int, season int) ([]model.PlayerGameStats, error) {
	rows, err := p.db.Query(
		"SELECT player_game_stats.player_id, player.name, game.id, player_game_stats.points, player_game_stats.assists, player_game_stats.rebounds, player_game_stats.steals, player_game_stats.blocks, player_game_stats.turnovers, player_game_stats.fouls, player_game_stats.minutes_played "+
			"FROM player_game_stats "+
			"JOIN player ON player_game_stats.player_id = player.id "+
			"JOIN game ON player_game_stats.game_id = game.id "+
			"JOIN seasons ON game.season_id = seasons.id "+
			"WHERE player_game_stats.team_id = $1 AND seasons.year = $2",
		teamID, season,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	var statsList []model.PlayerGameStats

	for rows.Next() {
		var stats model.PlayerGameStats
		err := rows.Scan(
			&stats.PlayerID,
			&stats.PlayerName,
			&stats.GameID,
			&stats.Points,
			&stats.Assists,
			&stats.Rebounds,
			&stats.Steals,
			&stats.Blocks,
			&stats.Turnovers,
			&stats.Fouls,
			&stats.MinutesPlayed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}
		statsList = append(statsList, stats)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return statsList, nil
}

// GetPlayerGamesBySeason implements PlayerRepository.
func (p *PlayerRepositoryStruct) GetPlayerGamesBySeason(playerID int, season int) ([]model.PlayerGameStats, error) {
	rows, err := p.db.Query(
		"SELECT player_game_stats.player_id, player.name, game.id, player_game_stats.points, player_game_stats.assists, player_game_stats.rebounds, player_game_stats.steals, player_game_stats.blocks, player_game_stats.turnovers, player_game_stats.fouls, player_game_stats.minutes_played "+
			"FROM player_game_stats "+
			"JOIN player ON player_game_stats.player_id = player.id "+
			"JOIN game ON player_game_stats.game_id = game.id "+
			"JOIN seasons ON game.season_id = seasons.id "+
			"WHERE player.id = $1 AND seasons.year = $2 "+
			"ORDER BY game.id ASC", // Ordering remains as is, adjust if needed
		playerID, season,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statsList []model.PlayerGameStats

	for rows.Next() {
		var stats model.PlayerGameStats
		err := rows.Scan(
			&stats.PlayerID,
			&stats.PlayerName,
			&stats.GameID,
			&stats.Points,
			&stats.Assists,
			&stats.Rebounds,
			&stats.Steals,
			&stats.Blocks,
			&stats.Turnovers,
			&stats.Fouls,
			&stats.MinutesPlayed,
		)
		if err != nil {
			return nil, err
		}
		statsList = append(statsList, stats)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statsList, nil
}

// LogPlayerGame implements PlayerRepository.
func (p *PlayerRepositoryStruct) LogPlayerGame(game model.PlayerGameStats) error {
	_, err := p.db.Exec(
		"INSERT INTO player_game_stats (player_id, game_id, points, assists, rebounds, steals, blocks, turnovers, fouls, minutes_played) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		game.PlayerID, game.GameID, game.Points, game.Assists, game.Rebounds, game.Steals, game.Blocks, game.Turnovers, game.Fouls, game.MinutesPlayed,
	)
	if err != nil {
		return err
	}
	return nil
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &PlayerRepositoryStruct{
		db: db,
	}
}
