package nats

import (
	"context"
	"xgo/main/src/event"

	"github.com/gookit/config/v2"
	natsClient "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Nats struct {
	producer *natsClient.Conn
}

func New() *Nats {
	conn, err := natsClient.Connect(config.String("producer.uri"), []natsClient.Option{}...)
	if err != nil {
		panic(err)
	}
	return &Nats{producer: conn}
}

func (n *Nats) Produce(ctx context.Context, topic string, ingestChan chan event.Event) {
	defer n.producer.Close()
	defer n.producer.Flush()
	for {
		select {
		case <-ctx.Done():
		case event := <-ingestChan:
			if err := n.producer.Publish(topic, event.Raw()); err != nil {
				zap.L().Error("error publish nats message", zap.Error(err))
			}
		}
	}
}
