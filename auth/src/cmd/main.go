// Copyright 2016 Google, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"xgo/auth/src/gen"
	grpchandler "xgo/auth/src/internal/handler/grpc"
	"xgo/auth/src/internal/handler/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger, _ := zap.NewProduction()
	origianlHook := zap.ReplaceGlobals(logger)
	defer origianlHook()

	initConfig()

	initConfig()

	repo := repository.New(config.String("auth.driver"))
	jwtPrivateKey, err := os.ReadFile(config.String("auth.jwt-private-key"))
	if err != nil {
		panic(err)
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(jwtPrivateKey))
	if err != nil {
		panic(err)
	}
	as := grpchandler.New(rsaKey, repo)
	srv := grpc.NewServer()

	reflection.Register(srv)
	gen.RegisterAuthServer(srv, as)
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.String("auth.port")))
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Serve(ln); err != nil {
		log.Fatalln("Failed to start the gRPC server", err)
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
