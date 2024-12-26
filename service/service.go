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
	GetPlayerSeasonAverages(ctx context.Context, request model.GetPlayerGameStatsRequest) (*model.PlayerSeasonAverage, error)
	GetTeamSeasonAverages(ctx context.Context, request model.GetTeamGameStatsRequest) (*model.TeamSeasoAverage, error)
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

	if playerId <= 0 {
		return errors.New("player ID must be a positive integer")
	}

	if request.GameId <= 0 {
		return errors.New("game ID must be a positive integer")
	}

	g, err := s.playerRepository.GetGame(request.GameId)
	if err != nil {
		return err
	}
	if g.Id == 0 {
		return errors.New("game not found")
	}

	p, err := s.playerRepository.GetPlayer(playerId)
	if err != nil {
		return err
	}
	if p.Id == 0 {
		return errors.New("player not found")
	}
	// If all validations pass, proceed with logging the game.
	// (Logging logic would go here, e.g., save the game stats to the database)
	playerGame := model.PlayerGameStats{
		PlayerID:      playerId,
		GameID:        request.GameId,
		Points:        request.Points,
		Assists:       request.Assists,
		Rebounds:      request.Rebounds,
		Steals:        request.Steals,
		Blocks:        request.Blocks,
		Turnovers:     request.Turnovers,
		Fouls:         request.Fouls,
		MinutesPlayed: request.MinutesPlayed,
	}
	err = s.playerRepository.LogPlayerGame(playerGame)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceStruct) GetPlayerSeasonAverages(ctx context.Context, req model.GetPlayerGameStatsRequest) (*model.PlayerSeasonAverage, error) {
	// Validate pagination parameters
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	// Limit the page size to 100
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Validate the season and player ID
	if req.SeasonYear <= 0 {
		return nil, errors.New("season must be a positive integer")
	}
	if req.PlayerID <= 0 {
		return nil, errors.New("player ID must be a positive integer")
	}

	// Get player data from the repository
	p, err := s.playerRepository.GetPlayer(req.PlayerID)
	if err != nil {
		return nil, err
	}
	if p.Id == 0 {
		return nil, errors.New("player not found")
	}

	// Get player stats by season
	playerStats, err := s.playerRepository.GetPlayerGamesBySeason(req.PlayerID, req.SeasonYear)
	if err != nil {
		return nil, err
	}
	totalGames := len(playerStats)
	if totalGames == 0 {
		return nil, nil
	}

	// Initialize the stats object with zero values
	stats := model.PlayerSeasonAverage{
		PlayerID:   p.Id,
		PlayerName: p.Name,
		Season:     req.SeasonYear,
	}

	// Accumulate the stats for all games
	for _, game := range playerStats {
		stats.PointsPerGame += float32(game.Points)
		stats.AssistsPerGame += float32(game.Assists)
		stats.ReboundsPerGame += float32(game.Rebounds)
		stats.StealsPerGame += float32(game.Steals)
		stats.BlocksPerGame += float32(game.Blocks)
		stats.TurnoversPerGame += float32(game.Turnovers)
		stats.FoulsPerGame += float32(game.Fouls)
		stats.MinutesPlayedPerGame += game.MinutesPlayed
	}

	// Calculate the averages
	if totalGames > 0 {
		stats.PointsPerGame /= float32(totalGames)
		stats.AssistsPerGame /= float32(totalGames)
		stats.ReboundsPerGame /= float32(totalGames)
		stats.StealsPerGame /= float32(totalGames)
		stats.BlocksPerGame /= float32(totalGames)
		stats.TurnoversPerGame /= float32(totalGames)
		stats.FoulsPerGame /= float32(totalGames)
		stats.MinutesPlayedPerGame /= float32(totalGames)
	}

	return &stats, nil
}

func (s *ServiceStruct) GetTeamSeasonAverages(ctx context.Context, req model.GetTeamGameStatsRequest) (*model.TeamSeasoAverage, error) {
	// Validate the pagination parameters
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Validate required fields
	if req.SeasonYear <= 0 {
		return nil, errors.New("season must be a positive integer")
	}
	if req.TeamID <= 0 {
		return nil, errors.New("team ID must be a positive integer")
	}

	// Retrieve the team from the repository
	team, err := s.playerRepository.GetTeam(req.TeamID)
	if err != nil {
		return nil, err
	}
	if team.Id == 0 {
		return nil, errors.New("team not found")
	}

	// Get team players' game stats for the season
	teamPlayersBySeason, err := s.playerRepository.GetTeamPlayersBySeason(req.TeamID, req.SeasonYear)
	if err != nil {
		return nil, err
	}
	totalGames := len(teamPlayersBySeason)
	if totalGames == 0 {
		return nil, nil
	}

	// Initialize the stats object
	stats := model.TeamSeasoAverage{}
	for _, teamPlayer := range teamPlayersBySeason {
		stats.PointsPerGame += float32(teamPlayer.Points)
		stats.AssistsPerGame += float32(teamPlayer.Assists)
		stats.ReboundsPerGame += float32(teamPlayer.Rebounds)
		stats.StealsPerGame += float32(teamPlayer.Steals)
		stats.BlocksPerGame += float32(teamPlayer.Blocks)
		stats.TurnoversPerGame += float32(teamPlayer.Turnovers)
		stats.FoulsPerGame += float32(teamPlayer.Fouls)
		stats.MinutesPlayedPerGame += float32(teamPlayer.MinutesPlayed)
	}

	// Calculate the averages, ensure no division by zero
	stats.PointsPerGame /= float32(totalGames)
	stats.AssistsPerGame /= float32(totalGames)
	stats.ReboundsPerGame /= float32(totalGames)
	stats.StealsPerGame /= float32(totalGames)
	stats.BlocksPerGame /= float32(totalGames)
	stats.TurnoversPerGame /= float32(totalGames)
	stats.FoulsPerGame /= float32(totalGames)
	stats.MinutesPlayedPerGame /= float32(totalGames)

	stats.TeamID = team.Id
	stats.TeamName = team.Name
	stats.Season = req.SeasonYear

	return &stats, nil
}
