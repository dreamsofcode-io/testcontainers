package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var ErrNoMessage = errors.New("no message on event bus")

type Message struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PubSub struct {
	conn *kafka.Conn
}

func New(conn *kafka.Conn) *PubSub {
	return &PubSub{
		conn: conn,
	}
}

func (p *PubSub) WriteMessage(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if _, err = p.conn.Write(data); err != nil {
		return fmt.Errorf("failed to write message %w", err)
	}

	return nil
}

func (p *PubSub) ReadMessage(ctx context.Context) (Message, error) {
	p.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	msg, err := p.conn.ReadMessage(1e6)
	if err != nil {
		return Message{}, fmt.Errorf("failed to read from connection: %w", err)
	}

	if msg.Value == nil {
		return Message{}, ErrNoMessage
	}

	var message Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return Message{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return message, nil
}
