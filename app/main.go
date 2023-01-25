package main

import (
	"fmt"
	"net"
	"time"

	"usermanager/app/config"
	"usermanager/app/infrastructure/db"
	notif "usermanager/app/infrastructure/notification"
	"usermanager/app/infrastructure/rabbit"
	repo "usermanager/app/infrastructure/repositories"
	"usermanager/app/services"
	h "usermanager/app/ui/grpcServers/health"
	u "usermanager/app/ui/grpcServers/user"

	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
)

func main() {
	// mock wait for the rmq containers to start properly
	time.Sleep(time.Second * 15)

	// load env variables
	config.Load()

	// run database migrations
	if err := db.MigrateDb(); err != nil {
		log.Fatal().Err(err).Msg("cannot run db migrations")
	}

	// start grpc server
	runGrpcServer()
}

func runGrpcServer() {
	port := fmt.Sprintf(":%v", config.EnvConfig.ServerPort)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to listen on port %v",
			config.EnvConfig.ServerPort)
	}

	// connect to db
	db, err := db.OpenDb()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// create repo, rmq wrapper, notification and user service
	userRepo := repo.NewUserRepo(db)
	rmq := rabbit.NewRMQ()
	notifService := notif.NewNotificationService(rmq)
	userService := services.NewUserService(userRepo, notifService)

	g := grpc.NewServer()

	// create and register user grpc server
	u.NewUserGrpcServer(g, userService)

	// create and register health grpc server
	h.NewHealthGrpcServer(g)

	if err := g.Serve(lis); err != nil {
		log.Fatal().Err(err).Msgf("Failed to serve gRPC server over port %v",
			config.EnvConfig.ServerPort)
	}
}
