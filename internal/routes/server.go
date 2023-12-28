package routes

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	socketio "github.com/vchitai/go-socket.io/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
	envProd        = "prod"
	certFile       = "/etc/letsencrypt/live/quizmaster.swedencentral.cloudapp.azure.com/fullchain.pem"
	keyFile        = "/etc/letsencrypt/live/quizmaster.swedencentral.cloudapp.azure.com/privkey.pem"
)

type Server struct {
	echo       *echo.Echo
	log        *zap.SugaredLogger
	db         *sqlx.DB
	rdb        *redis.Client
	esClient   *elasticsearch.Client
	ws         *socketio.Server
	tracer     trace.Tracer
	awsClient  *minio.Client
	rabbitConn *amqp.Connection
}

func NewServer(echo *echo.Echo, log *zap.SugaredLogger, db *sqlx.DB, rdb *redis.Client, esClient *elasticsearch.Client, ws *socketio.Server, tracer trace.Tracer, awsClient *minio.Client, rabbitConn *amqp.Connection) *Server {
	return &Server{echo: echo, log: log, db: db, rdb: rdb, esClient: esClient, ws: ws, tracer: tracer, awsClient: awsClient, rabbitConn: rabbitConn}
}

func (s *Server) Run() error {
	s.InitRoutes(s.echo)

	if os.Getenv("ENV") == envProd {
		go func() {
			s.log.Info("Server is listening on port 443")
			s.echo.Server.ReadTimeout = time.Second * 10
			s.echo.Server.WriteTimeout = time.Second * 10
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(":443", certFile, keyFile); err != nil {
				s.log.Fatalf("Error starting TLS Server: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer shutdown()

		s.log.Info("Server exiting...")
		return s.echo.Server.Shutdown(ctx)
	}

	go func() {
		if err := s.echo.Start(os.Getenv("PORT")); err != nil {
			s.log.Infof("Error Starting server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("Server exiting...")
	return s.echo.Server.Shutdown(ctx)
}
