package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/gorm"

	"github.com/rohit21755/gg_server.git/graph"
	"github.com/rohit21755/gg_server.git/ws"
)

// graphqlHandler wires resolvers, db, wsHub into a gqlgen HTTP handler
func graphqlHandler(db *gorm.DB, hub *ws.Hub) http.Handler {

	// Resolver struct contains all services (db, pubsub, etc.)
	resolver := &graph.Resolver{
		DB:    db,
		WsHub: hub,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	return srv
}

// graphqlPlayground provides an in-browser IDE
func graphqlPlayground() http.Handler {
	return playground.Handler("GraphQL Playground", "/graphql")
}

// mountGraphQL mounts all GraphQL endpoints
func mountGraphQL(mux *http.ServeMux, db *gorm.DB, hub *ws.Hub) {

	log.Println("GraphQL server mounted at /graphql")

	mux.Handle("/graphql", graphqlHandler(db, hub))
	mux.Handle("/playground", graphqlPlayground())
}
