//go:generate gqlgen -schema ./schema.graphql
package server

import (
	context "context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/ksuid"
	"github.com/vektah/gqlgen/handler"
)

type contextKey string

const (
	userContextKey = contextKey("user")
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
		messageCh:   make(chan Message, 1),
	}, nil
}

func (s *graphQLServer) Serve(route string, port int) error {
	http.Handle(
		route,
		authenticate(handler.GraphQL(MakeExecutableSchema(s))),
	)
	http.Handle("/playground", handler.Playground("GraphQL", route))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *graphQLServer) Mutation_postMessage(ctx context.Context, text string) (*Message, error) {
	user := ctx.Value(userContextKey).(string)
	if err := s.redisClient.LPush("users", user).Err(); err != nil {
		return nil, err
	}

	m := Message{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Text:      text,
		User:      user,
	}
	mj, _ := json.Marshal(m)
	if err := s.redisClient.LPush("messages", mj).Err(); err != nil {
		return nil, err
	}
	s.messageCh <- m

	return &m, nil
}

func (s *graphQLServer) Query_messages(ctx context.Context) ([]Message, error) {
	cmd := s.redisClient.LRange("messages", 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	messages := []Message{}
	for _, mj := range res {
		var m Message
		err = json.Unmarshal([]byte(mj), &m)
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *graphQLServer) Query_users(ctx context.Context) ([]string, error) {
	cmd := s.redisClient.LRange("users", 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *graphQLServer) Subscription_messagePosted(ctx context.Context) (<-chan Message, error) {
	return s.messageCh, nil
}
