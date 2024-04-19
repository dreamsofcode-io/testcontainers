package pubsub_test

import (
	"context"
	"log"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	kfka "github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/dreamsofcode-io/testcontainers/pubsub"
)

func TestPubSub(t *testing.T) {
	ctx := context.Background()

	container, err := kfka.RunContainer(ctx,
		kfka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
		testcontainers.WithEnv(map[string]string{
			"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE": "true",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("9093/tcp"),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := container.PortEndpoint(ctx, "9093/tcp", "")
	if err != nil {
		log.Fatalf("failed to get endpoint: %s", err)
	}

	conn, err := kafka.DialLeader(ctx, "tcp", endpoint, "test-topic", 0)
	if err != nil {
		log.Fatalf("failed to dial leader: %s", err)
	}

	t.Run("single message", func(t *testing.T) {
		ps := pubsub.New(conn)
		err = ps.WriteMessage(pubsub.Message{
			Title:       "Hello, world!",
			Description: "testcontainers are awesome",
		})

		assert.NoError(t, err)

		msg, err := ps.ReadMessage(ctx)
		assert.NoError(t, err)

		assert.Equal(t, "Hello, world!", msg.Title)
		assert.Equal(t, "testcontainers are awesome", msg.Description)
	})

	t.Run("multiple messages", func(t *testing.T) {
		ps := pubsub.New(conn)
		err = ps.WriteMessage(pubsub.Message{
			Title:       "Hello, world!",
			Description: "testcontainers are awesome",
		})

		assert.NoError(t, err)

		err = ps.WriteMessage(pubsub.Message{
			Title:       "Another one",
			Description: "golang is neat too",
		})

		assert.NoError(t, err)

		msg, err := ps.ReadMessage(ctx)
		assert.NoError(t, err)

		assert.Equal(t, "Hello, world!", msg.Title)
		assert.Equal(t, "testcontainers are awesome", msg.Description)

		msg, err = ps.ReadMessage(ctx)
		assert.NoError(t, err)

		assert.Equal(t, "Another one", msg.Title)
		assert.Equal(t, "golang is neat too", msg.Description)
	})
}
