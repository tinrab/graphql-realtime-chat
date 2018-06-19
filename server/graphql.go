//go:generate gqlgen -schema ./schema.graphql
package server

import (
	context "context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
	"github.com/vektah/gqlgen/handler"
)

type contextKey string

const (
	userContextKey = contextKey("user")
)

type graphQLServer struct {
	redisClient     *redis.Client
	messageChannels map[string]chan Message
	mutex           sync.Mutex
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
		redisClient:     client,
		messageChannels: map[string]chan Message{},
		mutex:           sync.Mutex{},
	}, nil
}

func (s *graphQLServer) Serve(route string, port int) error {
	mux := http.NewServeMux()
	mux.Handle(
		route,
		handler.GraphQL(MakeExecutableSchema(s),
			handler.WebsocketUpgrader(websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}),
		),
	)
	mux.Handle("/playground", handler.Playground("GraphQL", route))

	handler := cors.AllowAll().Handler(mux)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func (s *graphQLServer) Mutation_postMessage(ctx context.Context, user string, text string) (*Message, error) {
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

	s.mutex.Lock()
	for _, ch := range s.messageChannels {
		ch <- m
	}
	s.mutex.Unlock()

	return &m, nil
}

func (s *graphQLServer) Subscription_messagePosted(ctx context.Context, user string) (<-chan Message, error) {
	if err := s.redisClient.LPush("users", user).Err(); err != nil {
		return nil, err
	}

	messages := make(chan Message, 1)

	s.mutex.Lock()
	s.messageChannels[user] = messages
	s.mutex.Unlock()

	go func() {
		<-ctx.Done()
		s.mutex.Lock()
		delete(s.messageChannels, user)
		s.mutex.Unlock()
	}()

	return messages, nil
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
