package main

import (
	"log"
	"net/http"

	"github.com/tinrab/fluffy-kitten/graph"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	s := graph.NewGraphQLServer()
	http.Handle("/graphql", handler.GraphQL(graph.MakeExecutableSchema(s)))
	http.Handle("/playground", handler.Playground("App", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
