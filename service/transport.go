package service

import (
	"context"
	"log"
	"nba/model"
	"nba/pb"

	"go.uber.org/zap"
)

type GRPCServer struct {
	Logger *zap.SugaredLogger
	Svc    Service
	pb.UnimplementedPlayerGameServiceServer
}

func NewGRPCServer(logger *zap.SugaredLogger, svc Service) *GRPCServer {
	return &GRPCServer{Logger: logger, Svc: svc}
}

// Implement the LogPlayerGame method
func (t *GRPCServer) LogPlayerGame(ctx context.Context, req *pb.LogPlayerGameRequest) (*pb.LogGameResponse, error) {
	t.Logger.Info("Received LogPlayerGame request", req)
	request := model.LogPlayerGameRequest{
		PlayerId:      int(req.PlayerId),
		GameId:        int(req.GameId),
		Points:        int(req.Points),
		Rebounds:      int(req.Rebounds),
		Assists:       int(req.Assists),
		Steals:        int(req.Steals),
		Blocks:        int(req.Blocks),
		Turnovers:     int(req.Turnovers),
		Fouls:         int(req.Fouls),
		MinutesPlayed: req.MinutesPlayed,
	}

	err := t.Svc.LogPlayerGame(ctx, request.PlayerId, request)
	if err != nil {
		return nil, err
	}
	return &pb.LogGameResponse{Success: true}, nil
}

// Implement the GetPlayersGameStats method
func (t *GRPCServer) GetPlayerGameSeasonStats(ctx context.Context, request *pb.GetPlayerGameSeasonStatsRequest) (*pb.PlayerGameSeasonStatsResponse, error) {
	log.Println("Received GetPlayerGameSeasonStats request", request)
	player, err := t.Svc.GetPlayerSeasonAverages(ctx, model.GetPlayerGameStatsRequest{
		PlayerID:   int(request.PlayerId),
		SeasonYear: int(request.Season),
	})
	if err != nil {
		return nil, err
	}
	playerStats := pb.PlayerGameStat{
		Points:        int32(player.PointsPerGame),
		Assists:       int32(player.AssistsPerGame),
		Rebounds:      int32(player.ReboundsPerGame),
		Steals:        int32(player.StealsPerGame),
		Blocks:        int32(player.BlocksPerGame),
		Turnovers:     int32(player.TurnoversPerGame),
		Fouls:         int32(player.FoulsPerGame),
		MinutesPlayed: player.MinutesPlayedPerGame,
		PlayerId:      int32(player.PlayerID),
	}

	return &pb.PlayerGameSeasonStatsResponse{PlayerGameStats: &playerStats}, nil

}

func (t *GRPCServer) GetPlayer(ctx context.Context, request *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	log.Println("Received GetPlayerGameStats request", request)
	return &pb.GetPlayerResponse{Message: "cool", Success: true}, nil
}
