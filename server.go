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
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

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

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
