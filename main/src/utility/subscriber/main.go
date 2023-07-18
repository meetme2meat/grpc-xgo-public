package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	readConfig()
	fmt.Println("starting subscriber")
	conn, err := nats.Connect("nats://nats:4222", nil)
	if err != nil {
		panic(err)
	}
	reap(ctx, conn)
}

func readConfig() {
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

func reap(ctx context.Context, conn *nats.Conn) {
	sub, err := conn.SubscribeSync("event")
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := sub.NextMsg(2 * time.Second)
			if err != nil && errors.Is(err, nats.ErrTimeout) {
				continue
			} else if err != nil {
				panic(err)
			}

			fmt.Println("received event", string(msg.Data))
		}
	}
}
