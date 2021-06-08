package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/clshu/srv-go/dbmgm"
	"github.com/clshu/srv-go/graph/generated"
	"github.com/clshu/srv-go/graph/resolver"
	"github.com/clshu/srv-go/middleware/api"
	"github.com/clshu/srv-go/middleware/auth"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"
const defaultEndpoint = "tvu_graphql"

func main() {
	setUpEnv()

	err := dbmgm.Connect()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	endpoint := os.Getenv("GRAPHQL_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	gqlPath := fmt.Sprintf("/%s", endpoint)
	// log.Println(gqlPath)
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(auth.Middleware())
	r.Use(api.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))

	r.Handle("/tvu_playground", playground.Handler("GraphQL playground", gqlPath))
	r.Handle(gqlPath, srv)

	log.Printf("connect to http://localhost:%s/tvu_playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setUpEnv() {
	const dir string = "config/"
	var fname string
	gogo := os.Getenv("GOGO_ENV")

	switch gogo {
	case "dev":
		fname = dir + "dev.env"
	case "test":
		fname = dir + "test.env"
	default:
		// production environment
		// Do nothing. Let clound platform environment take over
		return
	}

	envMap, err := godotenv.Read(fname)
	if err != nil {
		fmt.Printf("Reading file %v failed. %v", fname, err.Error())
		return
	}

	for key, value := range envMap {
		os.Setenv(key, value)
		// fmt.Printf("%v=%v\n", key, os.Getenv(key))
	}

}
