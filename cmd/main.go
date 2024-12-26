package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq" // or "github.com/jackc/pgx/v5/stdlib"

	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"nba/pb"
	p "nba/postgres"
	"nba/service"

	"go.uber.org/zap"
)

// Helper function to handle errors
func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	checkError(err, "Failed to create logger")
	defer logger.Sync()

	logger.Info("Starting the NBA service")

	// Fetch database connection details from environment variables (with defaults)
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	port := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "nba")

	// Step 2: After ensuring the database exists, connect to it using the correct dbname
	connStr := fmt.Sprintf("user=%s password=%s host=db port=%s dbname=%s sslmode=disable", user, password, port, dbName)
	db, err := sql.Open("postgres", connStr)
	checkError(err, "Failed to connect to the specific database")
	defer db.Close()

	// Step 4: Initialize the database
	if err := initDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a new player repository
	playerRepository := p.NewPlayerRepository(db)

	// Create a new service
	svc := service.NewService(logger.Sugar(), playerRepository)

	logger.Info("Starting the NBA service")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register the PlayerGameService with the gRPC server
	pb.RegisterPlayerGameServiceServer(grpcServer, &service.GRPCServer{Logger: logger.Sugar(), Svc: svc})

	// Start a gRPC server on port 50051
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Create HTTP Gateway
	mux := runtime.NewServeMux()

	// Register the HTTP Gateway for PlayerGameService
	// You need to specify a gRPC client connection to the server
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	checkError(err, "Failed to dial gRPC server")
	defer grpcConn.Close()

	// Register the gRPC service to the HTTP Gateway mux
	err = pb.RegisterPlayerGameServiceHandlerFromEndpoint(context.Background(), mux, "localhost:50051", []grpc.DialOption{grpc.WithInsecure()})

	// Start HTTP server on port 8080
	http.Handle("/", mux)
	log.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}
func initDatabase(db *sql.DB) error {
	// SQL statements to create the tables
	createPlayerTable := `
	CREATE TABLE IF NOT EXISTS player (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		current_team_id INT NOT NULL,
		CONSTRAINT fk_team FOREIGN KEY (current_team_id) REFERENCES team (id) ON DELETE CASCADE
	);`

	createPlayerGameStatTable := `
	CREATE TABLE IF NOT EXISTS player_game_stats (
		player_id INT NOT NULL,
		game_id INT NOT NULL,
		points INT,
		assists INT,
		rebounds INT,
		steals INT,
		blocks INT,
		turnovers INT,
		fouls INT,
		minutes_played FLOAT,
		PRIMARY KEY (player_id, game_id),
		CONSTRAINT fk_player FOREIGN KEY (player_id) REFERENCES player (id) ON DELETE CASCADE,
		CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES game (id) ON DELETE CASCADE
	);`

	createGameTable := `
	CREATE TABLE IF NOT EXISTS game (
		id SERIAL PRIMARY KEY,
		date DATE NOT NULL,
		season INT NOT NULL,
		team_a_id INT NOT NULL,
		team_b_id INT NOT NULL
	);`

	createTeamTable := `
	CREATE TABLE IF NOT EXISTS team (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL
	);`

	// Execute the SQL statements
	tables := []string{
		createTeamTable,
		createGameTable,
		createPlayerTable,
		createPlayerGameStatTable,
	}

	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v\nquery: %s", err, query)
		}
	}

	// Insert teams
	teams := []string{
		"INSERT INTO team (name) VALUES ('Team A'), ('Team B'), ('Team C');",
	}

	for _, query := range teams {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to insert teams: %v", err)
		}
	}

	// Insert games
	_, err := db.Exec(`
		INSERT INTO game (date, season, team_a_id, team_b_id)
		VALUES
			('2024-01-01', 2024, 1, 2),
			('2024-01-02', 2024, 2, 3),
			('2024-01-03', 2024, 3, 1);
	`)
	if err != nil {
		return fmt.Errorf("failed to insert games: %v", err)
	}

	// Insert players
	_, err = db.Exec(`
		INSERT INTO player (name, current_team_id)
		VALUES
			('Player 1', 1),
			('Player 2', 2);
	`)
	if err != nil {
		return fmt.Errorf("failed to insert player: %v", err)
	}

	return nil
}
