package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"xgo/main/src/event"
	"xgo/main/src/gen"

	"xgo/main/src/internal/controller/company"
	grpchandler "xgo/main/src/internal/handler/grpc"
	"xgo/main/src/internal/producer"
	"xgo/main/src/internal/repository"

	"github.com/golang-jwt/jwt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ingestChan = make(chan event.Event)

func main() {
	logger, _ := zap.NewProduction()
	origianlHook := zap.ReplaceGlobals(logger)
	defer origianlHook()

	initConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := config.Int("main.port")

	defer close(ingestChan)
	driver := config.String("main.driver")

	publickey := getJWTPublicKey()
	repo := repository.New(driver, ingestChan)
	ctrl := company.New(repo)

	h := grpchandler.New(ctrl, publickey)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	// tlsCert := config.String("main.tlscert")
	// tlsKey := config.String("main.tlskey")
	// creds = tls.GetTransportCredentials(tlsCert, tlsKey)

	zap.L().Info("starting grpc server 0.0.0.0", zap.Any("port", port))
	notifier := make(chan os.Signal, 2)
	signal.Notify(notifier, syscall.SIGHUP, syscall.SIGTERM)
	go monitorSignal(ctx, notifier, cancel)

	// Start the ingestion worker
	producerWorker := producer.New(config.String("producer.name"))
	go startIngester(ctx, producerWorker)

	// srv := grpc.NewServer(grpc.Creds(creds))
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterCompanyServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		log.Fatalln("Failed to start the gRPC server", err)
	}
}

func startIngester(ctx context.Context, handler producer.Producer) {
	handler.Produce(ctx, config.String("producer.topic"), ingestChan)
}

func monitorSignal(ctx context.Context, notifier chan os.Signal, cancel context.CancelFunc) {
	switch <-notifier {
	case syscall.SIGINT, syscall.SIGTERM:
		zap.L().Info("received SIGTERM | SIGINT")
		signal.Stop(notifier)
		cancel()
	}
}

func initConfig() {
	config.AddDriver(ini.Driver)
	config.WithOptions(
		config.ParseEnv,
		config.WithHookFunc(func(event string, c *config.Config) {
			if event == config.OnReloadData {
				zap.L().Debug("config is reloaded")
			}
		}),
	)
	zap.L().Debug("LoadFiles...")
	err := config.LoadFiles("xgo.ini")
	if err != nil {
		panic(err)
	}
}

func getJWTPublicKey() *rsa.PublicKey {
	data, err := os.ReadFile(config.String("main.jwt-public-key"))
	if err != nil {
		log.Panicf("error reading the jwt public key: %v", err)
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		log.Panicf("Error parsing the jwt public key: %s", err)
	}

	return publickey
}
