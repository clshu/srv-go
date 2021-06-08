package gql

type GQLInfo struct {
	OperationName string
	Query         string
}

var CreateUser = GQLInfo{
	OperationName: "CreateUser",
	Query: string(`
  mutation CreateUser($data: CreateUserInput!) {
    createUser(data: $data) {
        id
        firstName
        lastName
        email
    }
  }
`),
}
