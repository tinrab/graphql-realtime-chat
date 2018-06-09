//go:generate gqlgen -schema ./schema.graphql
package server

import (
	context "context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/ksuid"
	"github.com/vektah/gqlgen/handler"
)

type graphQLServer struct {
	redisClient *redis.Client
	messageCh   chan Message
}

func NewGraphQLServer(redisURL string) (*graphQLServer, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &graphQLServer{
		redisClient: client,
		messageCh:   make(chan Message, 16),
	}, nil
}

func (s *graphQLServer) Serve(route string, port int) error {
	http.Handle(route, authenticate(handler.GraphQL(MakeExecutableSchema(s))))
	http.Handle("/__playground", handler.Playground("GraphQL", route))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *graphQLServer) Mutation_postMessage(ctx context.Context, text string) (*Message, error) {
	user := ctx.Value(userContextKey).(string)

	m := &Message{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Text:      text,
		User:      user,
	}

	s.messageCh <- *m

	return m, nil
}

func (s *graphQLServer) Query_users(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (s *graphQLServer) Subscription_messagePosted(ctx context.Context) (<-chan Message, error) {
	return s.messageCh, nil
}

func (s *graphQLServer) Message_user(ctx context.Context, m *Message) (*string, error) {
	return nil, nil
}
