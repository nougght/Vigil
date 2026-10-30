package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/backend/pkg/natsclient"
	"github.com/nougght/monitoring-system/backend/services/alerts/internal/config"
	"github.com/nougght/monitoring-system/backend/services/alerts/internal/service/event"
	"google.golang.org/grpc"
)

type App struct {
	Config *config.Config
	DB     *pgxpool.Pool

	//Repositories *repository.Repositories
	//
	//Services *service.Services

	HTTPServer       *http.Server
	AlertsGRPCServer *grpc.Server
}

func New(ctx context.Context, cfg *config.Config) *App {
	//db, err := storage.ConnectAndMigrate(ctx, cfg.Postgres, migrationFS)
	//if err != nil {
	//	log.Println("failed to connect to database, retry after 500ms")
	//	select {
	//	case <-time.After(time.Millisecond * 1000):
	//		db, err = storage.ConnectAndMigrate(ctx, cfg.Postgres, migrationFS)
	//		if err != nil {
	//			log.Panicf("failed to connect to database: %s", err.Error())
	//		}
	//	case <-ctx.Done():
	//		log.Panicf("%s", ctx.Err().Error())
	//	}

	//}
	//services := service.New(service.ServicesOptions{
	//	Config:       cfg,
	//	Repositories: repository.New(db),
	//	Transactor:   storage.NewTransactor(db),
	//})

	r := gin.New()
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.ServerPort),
		Handler: r.Handler(),
	}

	//alertsServer := grpc.NewServer(
	//	grpc.ConnectionTimeout(5*time.Second),
	//	grpc.StreamInterceptor(agent_grpc.BidirectionalServerInterceptor),
	//	grpc.Creds(
	//		credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{*cert},
	//			ClientAuth: tls.RequireAndVerifyClientCert, //  mTLS
	//			ClientCAs:  rootCA,
	//			MinVersion: tls.VersionTLS12,
	//		}),
	//	))

	//agentService := agent_grpc.NewAgentService(services.AgentInteractionService())
	//agentService.Register(agentServer)

	app := &App{
		Config: cfg,
		//DB:                   db,
		//Repositories:         repository.New(db),
		//Services:             services,
		//Bus:                  bus,
		HTTPServer: &httpServer,
		//AlertsGRPCServer:     alertsServer,
	}

	//go func() {
	//	<-ctx.Done()
	//	log.Println("broadcast shutdown")
	//	agentService.BroadcastShutdown()
	//}()
	return app
}

func (a *App) Run(ctx context.Context) error {

	natsClient := natsclient.NewClient(a.Config.NatsURL)
	eventService := event.NewEventService(natsClient.Durable())
	err := eventService.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed start event service: %w", err)
	}

	httpErrChan := make(chan error)
	go func() {
		log.Println("starting HTTP server")
		if err := a.HTTPServer.ListenAndServe(); err != nil {
			httpErrChan <- err
		}
	}()

	//l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.GRPC.ServerPort))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//grpcErrChan := make(chan error)
	//go func() {
	//	log.Printf("starting gRPC server :%d", a.Config.GRPC.ServerPort)
	//	if err := a.AlertsGRPCServer.Serve(l); err != nil {
	//		grpcErrChan <- err
	//	}
	//}()

	select {
	//case err := <-grpcErrChan:
	//	log.Fatalf("failed to serve gRPC  server: %v", err)
	case err := <-httpErrChan:
		log.Fatalf("failed to serve HTTP: %v", err)
	case <-ctx.Done():
		log.Println("root ctx done start shutdown")
	}

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	buf := make([]byte, 1<<20)
	// 	n := runtime.Stack(buf, true)
	// 	log.Printf("=== goroutines ===\n%s", buf[:n])
	// }()

	log.Println("start alerts grpc stop")
	a.AlertsGRPCServer.GracefulStop()
	log.Println("start http shutdwon")
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	err = a.HTTPServer.Shutdown(shutdownCtx)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to shutdown HTTP server: %v", err)
	}

	return nil
}
