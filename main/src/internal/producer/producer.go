package producer

import (
	"context"
	"xgo/main/src/event"
	"xgo/main/src/internal/producer/nats"
)

type Producer interface {
	Produce(context.Context, string, chan event.Event)
}

func New(name string) Producer {
	switch name {
	// Kafka build is causing some problem on alphine
	// case "kafka":
	// 	return kafka.New()
	case "nats":
		return nats.New()
	default:
		panic("not a known producer")
	}
}
