package service

import (
	"context"
	"errors"
	"nba/model"
	"nba/postgres"

	"go.uber.org/zap"
)

type Service interface {
	LogPlayerGame(ctx context.Context, playerId int, request model.LogPlayerGameRequest) error
	GetPlayerSeasonAverages(ctx context.Context, request model.GetPlayersGameStatsRequest) ([]model.PlayerSeasonAverage, int64, error)
	GetTeamSeasonStats(ctx context.Context, request model.GetTeamsGameStatsRequest) ([]model.TeamSeasoAverage, error)
}

type ServiceStruct struct {
	logger           *zap.SugaredLogger
	playerRepository postgres.PlayerRepository
}

func NewService(logger *zap.SugaredLogger, playerRepository postgres.PlayerRepository) Service {
	return &ServiceStruct{
		logger:           logger,
		playerRepository: playerRepository,
	}
}

// LogPlayerGame implements Service.
func (s *ServiceStruct) LogPlayerGame(ctx context.Context, playerId int, request model.LogPlayerGameRequest) error {
	// Validate Points, Rebounds, Assists, Steals, Blocks, Turnovers (must be positive integers)
	if request.Points <= 0 {
		return errors.New("points must be a positive integer")
	}
	if request.Rebounds <= 0 {
		return errors.New("rebounds must be a positive integer")
	}
	if request.Assists <= 0 {
		return errors.New("assists must be a positive integer")
	}
	if request.Steals <= 0 {
		return errors.New("steals must be a positive integer")
	}
	if request.Blocks <= 0 {
		return errors.New("blocks must be a positive integer")
	}
	if request.Turnovers <= 0 {
		return errors.New("turnovers must be a positive integer")
	}

	// Validate Fouls (must be an integer between 0 and 6)
	if request.Fouls < 0 || request.Fouls > 6 {
		return errors.New("fouls must be an integer between 0 and 6")
	}

	// Validate MinutesPlayed (must be a float between 0 and 48.0)
	if request.MinutesPlayed < 0 || request.MinutesPlayed > 48.0 {
		return errors.New("minutes played must be a float between 0 and 48.0")
	}

	// If all validations pass, proceed with logging the game.
	// (Logging logic would go here, e.g., save the game stats to the database)
	playerGame := model.PlayerGameStat{
		PlayerID:      playerId,
		Points:        request.Points,
		Assists:       request.Assists,
		Rebounds:      request.Rebounds,
		Steals:        request.Steals,
		Blocks:        request.Blocks,
		Turnovers:     request.Turnovers,
		Fouls:         request.Fouls,
		MinutesPlayed: request.MinutesPlayed,
	}
	err := s.playerRepository.LogPlayerGame(playerGame)
	if err != nil {
		return err
	}
	return nil

}

func (s *ServiceStruct) GetPlayerSeasonAverages(ctx context.Context, req model.GetPlayersGameStatsRequest) ([]model.PlayerSeasonAverage, int64, error) {
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	// limit the page size to 100
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	var players []model.Player
	players, count, err := s.playerRepository.GetPlayers(req.PageNumber, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []model.PlayerSeasonAverage
	for _, player := range players {
		totalGames := len(player.Games)
		if totalGames == 0 {
			continue
		}

		var stats model.PlayerGameStat
		for _, game := range player.Games {
			stats.Points += game.Points
			stats.Assists += game.Assists
			stats.Rebounds += game.Rebounds
			stats.Steals += game.Steals
			stats.Blocks += game.Blocks
			stats.Turnovers += game.Turnovers
			stats.Fouls += game.Fouls
			stats.MinutesPlayed += game.MinutesPlayed
		}

		results = append(results, model.PlayerSeasonAverage{
			PlayerID:             player.Id,
			TeamID:               player.TeamID, // Fill with team info if needed
			PointsPerGame:        float32(stats.Points) / float32(totalGames),
			AssistsPerGame:       float32(stats.Assists) / float32(totalGames),
			ReboundsPerGame:      float32(stats.Rebounds) / float32(totalGames),
			StealsPerGame:        float32(stats.Steals) / float32(totalGames),
			BlocksPerGame:        float32(stats.Blocks) / float32(totalGames),
			TurnoversPerGame:     float32(stats.Turnovers) / float32(totalGames),
			FoulsPerGame:         float32(stats.Fouls) / float32(totalGames),
			MinutesPlayedPerGame: stats.MinutesPlayed / float32(totalGames),
		})
	}

	return results, count, nil
}

func (s *ServiceStruct) GetTeamSeasonStats(ctx context.Context, req model.GetTeamsGameStatsRequest) ([]model.TeamSeasoAverage, error) {
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	teams, err := s.playerRepository.GetTeams()
	if err != nil {
		return nil, err
	}
	if len(teams) == 0 {
		return nil, nil
	}
	var results []model.TeamSeasoAverage
	for _, team := range teams {
		players, err := s.playerRepository.GetTeamPlayers(team.Id)
		if err != nil {
			return nil, err
		}

		var stats model.PlayerGameStat
		totalGames := 0
		for _, player := range players {
			totalGames += len(player.Games)
			for _, game := range player.Games {
				stats.Points += game.Points
				stats.Assists += game.Assists
				stats.Rebounds += game.Rebounds
				stats.Steals += game.Steals
				stats.Blocks += game.Blocks
				stats.Turnovers += game.Turnovers
				stats.Fouls += game.Fouls
				stats.MinutesPlayed += game.MinutesPlayed
			}
		}
		// Check if there are no games played by any player in the team, then skip the team
		if totalGames == 0 {
			continue
		}
		results = append(results, model.TeamSeasoAverage{
			TeamName:             team.Name,
			TeamID:               team.Id,
			PointsPerGame:        float32(stats.Points) / float32(totalGames),
			AssistsPerGame:       float32(stats.Assists) / float32(totalGames),
			ReboundsPerGame:      float32(stats.Rebounds) / float32(totalGames),
			StealsPerGame:        float32(stats.Steals) / float32(totalGames),
			BlocksPerGame:        float32(stats.Blocks) / float32(totalGames),
			TurnoversPerGame:     float32(stats.Turnovers) / float32(totalGames),
			FoulsPerGame:         float32(stats.Fouls) / float32(totalGames),
			MinutesPlayedPerGame: stats.MinutesPlayed / float32(totalGames),
		})
	}
	return results, nil

}
