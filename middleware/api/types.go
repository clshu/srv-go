package api

import (
	"github.com/clshu/srv-go/middleware/api/gql"
)

type GQLBody struct {
	OperationName string `json:"operationName,omitempty"`
	// Variables     map[string]interface{} `json:"variables,omitempty"`
	Variables interface{} `json:"variables,omitempty"`
	Query     string      `json:"query"`
}

type GQLB struct {
	OperationName string                 `json:"operaiotnName,omitempty"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// type VaraibalesType interface {
// 	DecodeInput(body io.ReadCloser)
// }

// REST API types
// Defined by Method + Path
// type APIType string

const (
	CreateUser string = "POST/app/user/create"
)

var APIFinder = map[string]gql.GQLInfo{
	CreateUser: gql.CreateUser,
}
