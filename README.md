# srv-go

## gqlge + mgm

## Steps to create the project

### go get github.com/99designs/gqlgen

### go get github.com/vektah/gqlparser/v2@v2.1.0

### go run github.com/99designs/gqlgen init

## Steps to change the project

### Modify gqlgen yml

#### resolver:

layout: follow-schema
dir: graph/resolver
package: resolver

### Copy graph/schema.graphqls to graph/user.graphqls and graph/todeo.graphqls

### Remove graph/schema.graphqls

### Remove graph/resolver.go

### Remove graph/schema.resolver.go

### go run github.com/99designs/gqlgen generate

## Rerun generate

### go get github.com/99designs/gqlgen/internal/imports@v0.13.0

### go get github.com/99designs/gqlgen/cmd@v0.13.0

### go get github.com/99designs/gqlgen/internal/code@v0.13.0
