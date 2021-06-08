package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/clshu/srv-go/graph/model"
)

// The middleware to convert a Rest API to a
// graphql call

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// var buf = GQLB{}
			// json.NewDecoder(r.Body).Decode(&buf)
			// log.Printf("%+v\n", buf)
			// log.Printf("%v\n", r.URL.Path)
			// log.Printf("%v\n", r.URL.Query())
			apiType := strings.Join([]string{r.Method, r.URL.Path}, "")
			// log.Printf("apiType: %s\n", apiType)
			gqlInfo := APIFinder[apiType]
			// log.Printf("gqlStr: %s\n", gqlInfo)
			if len(gqlInfo.Query) == 0 && len(gqlInfo.OperationName) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			// var input = &model.CreateUserInput{}
			// json.NewDecoder(r.Body).Decode(input)
			// log.Printf("%+v\n", r.Body)
			transformRequest(w, r, apiType)
			// if len(APIFinder[APIType(apiType)]) == 0 {
			// 	// key not found.
			// 	// Not a legit Rest API path
			// 	// Skip this middleware and continute to the next
			// 	next.ServeHTTP(w, r)
			// }

			// log.Printf("%s\n", PathFinder[key])

			// input := &model.CreateUserInput{}
			// err := json.NewDecoder(r.Body).Decode(input)
			// if err != nil {
			// 	returnError(&w, http.StatusBadRequest, err.Error())
			// 	return
			// }
			// log.Printf("input: %+v\n", input)
			next.ServeHTTP(w, r)

		})
	}
}

func transformRequest(w http.ResponseWriter, r *http.Request, apiType string) {

	opName := APIFinder[apiType].OperationName
	query := APIFinder[apiType].Query
	newBody := &GQLBody{OperationName: opName, Query: query}
	r.Method = "POST"
	r.URL.Path = "/tvu_graphql"

	switch apiType {
	case CreateUser:
		input := &model.CreateUserInput{}
		json.NewDecoder(r.Body).Decode(input)
		newBody.Variables = (*input)
	}
	jbody, _ := json.Marshal(newBody)
	// var buf bytes.Buffer
	// err := json.NewEncoder(&buf).Encode(in)
	// r.Body = ioutil.NopCloser(bytes.NewReader(jbody))
	r.Body = ioutil.NopCloser(strings.NewReader(string(jbody)))

	// log.Printf("%s\n", string(jbody))
	// jgqlBody := &GQLBody{}
	// json.Unmarshal(jbody, jgqlBody)
	// log.Printf("%+v\n", jgqlBody)

	// gqlBody := &GQLB{}
	// json.NewDecoder(r.Body).Decode(gqlBody)
	// log.Printf("%+v\n", gqlBody)
}

// func setVariables(r *http.Request, newBody *GQLBody, input interface{}) {
// 	json.NewDecoder(r.Body).Decode(input)
// 	newBody.Variables = (*input)
// }

// func replaceBody(w http.ResponseWriter, r *http.Request, input interface{}, newBody *GQLBody) {
// 	json.NewDecoder(r.Body).Decode(input)
// 	log.Printf("input:\n%+v\n", input)
// 	newBody.Variables = (*input)
// 	log.Printf("new body:\n%+v\n", newBody)
// 	jbody, _ := json.Marshal(newBody)
// 	// r.Body = ioutil.NopCloser(io.Reader(bytes.NewReader(jbody)))
// 	r.Body = ioutil.NopCloser(bytes.NewReader(jbody))

// }

// func DecodeInput(body io.ReadCloser, variables *VaraibalesType) {

// }

// type GraphqlReturn struct {
// 	Errors []string    `json:"errors"`
// 	Data   interface{} `json:"data"`
// }

// func returnError(w *http.ResponseWriter, status int, msg string) {
// 	(*w).WriteHeader(status)
// 	ret := &GraphqlReturn{Errors: []string{msg}, Data: nil}
// 	_ = json.NewEncoder(*w).Encode(ret)
// }
