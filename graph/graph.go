//go:generate gqlgen -schema ./schema.graphql
package graph

import (
	context "context"
	"time"
)

type graphQLServer struct {
}

func NewGraphQLServer() *graphQLServer {
	return &graphQLServer{}
}

func (s *graphQLServer) Mutation_postMessage(ctx context.Context, text string) (*Message, error) {
	return nil, nil
}

func (s *graphQLServer) Query_messages(ctx context.Context, skip *int, take *int) ([]Message, error) {
	return nil, nil
}

func (s *graphQLServer) Subscription_messagePosted(ctx context.Context) (<-chan Message, error) {
	result := make(chan Message)
	go func() {
		for {
			result <- Message{
				Text: "Hi",
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return result, nil
}

func (s *graphQLServer) Message_user(ctx context.Context, m *Message) (*User, error) {
	return nil, nil
}
