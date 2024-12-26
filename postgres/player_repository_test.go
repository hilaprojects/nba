package postgres_test

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"nba/model"
// 	db "nba/postgres"
// 	"os"
// 	"testing"
// 	"time"

// 	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
// 	"go.uber.org/mock/gomock"
// 	postgres "gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var container *postgresContainer.PostgresContainer
// var pgDB *gorm.DB // Global DB variable

// type PlayerRepositoryTestSuite struct {
// 	pgDB             *gorm.DB
// 	ctrl             *gomock.Controller
// 	playerRepository db.PlayerRepository
// }

// // TestMain initializes the PostgreSQL container and database connection
// func TestMain(m *testing.M) {
// 	// Define the PostgreSQL container configuration
// 	ctx := context.Background()
// 	var err error
// 	container, err = postgresContainer.Run(
// 		ctx,

// 		"postgres:latest",                              // Specify the Docker image to use
// 		postgresContainer.WithDatabase("nba"),          // Database name
// 		postgresContainer.WithUsername("testuser"),     // Database username
// 		postgresContainer.WithPassword("testpassword"), // Database password

// 	)
// 	if err != nil {
// 		log.Fatalf("Could not start PostgreSQL container: %s", err)
// 	}

// 	// Retry connection for 10 seconds (1 second delay between attempts)
// 	var port int
// 	for i := 0; i < 10; i++ {
// 		mappedPort, err := container.MappedPort(ctx, "5432")
// 		if err == nil {
// 			port = mappedPort.Int()
// 			break // Exit the loop if the port is successfully retrieved
// 		}
// 		log.Printf("Failed to get container port (attempt %d/10): %v", i+1, err)
// 		time.Sleep(1 * time.Second) // Wait before retrying
// 	}
// 	if port == 0 {
// 		log.Fatalf("Failed to get container port after multiple attempts")
// 	}

// 	// Connect to the database using GORM
// 	dsn := fmt.Sprintf("host=localhost port=%d user=testuser password=testpassword dbname=nba sslmode=disable", port)
// 	for i := 0; i < 10; i++ {
// 		pgDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 		if err == nil {
// 			log.Println("Connected to PostgreSQL successfully")
// 			break // Exit the loop if the connection is successful
// 		}
// 		log.Printf("Failed to connect to PostgreSQL (attempt %d/10): %v", i+1, err)
// 		time.Sleep(1 * time.Second) // Wait before retrying
// 	}

// 	if pgDB == nil {
// 		log.Fatalf("Failed to connect to PostgreSQL after multiple attempts")
// 	}

// 	// Inside your TestMain or a similar setup function:
// 	err = pgDB.AutoMigrate(&model.PlayerGameStat{})
// 	if err != nil {
// 		log.Fatalf("Error migrating database: %v", err)
// 	}
// 	// Run the tests
// 	code := m.Run()

// 	// Clean up the container after tests
// 	if err := container.Terminate(ctx); err != nil {
// 		log.Fatalf("Failed to terminate container: %v", err)
// 	}

// 	// Exit with the test result code
// 	os.Exit(code)
// }

// // getPlayerTestSuite sets up the test suite with the database connection
// func getPlayerTestSuite(t *testing.T) *PlayerRepositoryTestSuite {
// 	ctrl := gomock.NewController(t)

// 	// Create the player repository
// 	playerRepository := db.NewPlayerRepository(pgDB)

// 	// Return the test suite
// 	return &PlayerRepositoryTestSuite{
// 		pgDB:             pgDB,
// 		ctrl:             ctrl,
// 		playerRepository: playerRepository,
// 	}
// }

// // TestService_LogPlayerGame tests the player game logging functionality
// func TestService_LogPlayerGame(t *testing.T) {
// 	// Setup test suite
// 	tb := getPlayerTestSuite(t)

// 	// Save a team to the database
// 	tb.playerRepository.SaveTeam(model.Team{
// 		Id:   1,
// 		Name: "Lakers",
// 	})

// 	// Save a player to the database
// 	tb.playerRepository.SavePlayer(model.Player{
// 		Id:            1,
// 		Name:          "John",
// 		CurrentTeamID: 1,
// 	})
// 	// Perform the action: logging a player's game
// 	tb.playerRepository.LogPlayerGame(model.PlayerGameStat{
// 		PlayerID:      1,
// 		Points:        10,
// 		Assists:       5,
// 		Rebounds:      3,
// 		Steals:        2,
// 		Blocks:        1,
// 		Turnovers:     4,
// 		Fouls:         3,
// 		MinutesPlayed: 30.0,
// 	})

// 	db := tb.pgDB
// 	var playerGameInfo model.PlayerGameStat
// 	db.Exec("DELETE FROM players")
// 	if playerGameInfo.PlayerID != 1 {
// 		t.Errorf("Player ID does not match")
// 	}

// 	db.Exec("SELECT * FROM player_game_stats").First(&playerGameInfo)

// 	// Check if the game was logged (you can add an actual assertion here)
// 	fmt.Println("Player game logged successfully")
// }
