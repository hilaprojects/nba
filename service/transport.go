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
func (t *GRPCServer) GetPlayersGameStats(ctx context.Context, request *pb.GetPlayersGameStatsRequest) (*pb.PlayersGameStatsResponse, error) {
	log.Println("Received GetPlayersGameStats request", request)
	p, c, err := t.Svc.GetPlayerSeasonAverages(ctx, model.GetPlayersGameStatsRequest{
		PageSize:   int(request.PageSize),
		PageNumber: int(request.PageNumber),
	})
	if err != nil {
		return nil, err
	}
	var playerStats []*pb.PlayerGameStat
	if p != nil {
		for _, player := range p {
			playerStats = append(playerStats, &pb.PlayerGameStat{
				Points:        int32(player.PointsPerGame),
				Assists:       int32(player.AssistsPerGame),
				Rebounds:      int32(player.ReboundsPerGame),
				Steals:        int32(player.StealsPerGame),
				Blocks:        int32(player.BlocksPerGame),
				Turnovers:     int32(player.TurnoversPerGame),
				Fouls:         int32(player.FoulsPerGame),
				MinutesPlayed: player.MinutesPlayedPerGame,
				PlayerId:      int32(player.PlayerID),
			})
		}
	}
	return &pb.PlayersGameStatsResponse{PlayerGameStats: playerStats, TotalCount: c}, nil

}
