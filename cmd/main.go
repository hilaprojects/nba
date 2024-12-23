package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"nba/pb"
	p "nba/postgres"
	"nba/service"

	"go.uber.org/zap"
)

var (
	user     = os.Getenv("PSQL_DB_USER")
	password = os.Getenv("PSQL_DB_PASSWORD")
	port     = os.Getenv("PSQL_DB_PORT")
	host     = os.Getenv("PSQL_DB_HOST")
	dbname   = os.Getenv("PSQL_DB_NAME")
	sslmode  = os.Getenv("PSQL_DB_SSL_MODE")
)

func main() {

	// Create a production logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Flush any buffered log entries

	logger.Info("Starting the NBA service")

	// Create a new database connection
	db, err := ConnectPostgres(user, password, port, host, dbname, sslmode)
	if err != nil {
		logger.Error("failed to connect to the database", zap.Error(err))
		return
	}
	// Create a new player repository
	playerRepository := p.NewPlayerRepository(db)

	// Create a new service
	svc := service.NewService(logger.Sugar(), playerRepository)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register the service with the gRPC server
	pb.RegisterPlayerGameServiceServer(grpcServer, &service.GRPCServer{Logger: logger.Sugar(), Svc: svc})

	log.Println("Server is listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

// ConnectPostgres establishes a GORM connection to a PostgreSQL database.
func ConnectPostgres(user, password, port, host, dbname, sslmode string) (*gorm.DB, error) {
	// Build the connection string
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	// Initialize GORM with the PostgreSQL driver
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	log.Println("Successfully connected to the database using GORM")
	return db, nil
}
