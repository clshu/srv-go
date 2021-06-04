package utils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/clshu/go-mgm/utils"
)

func TestCreateToken(t *testing.T) {
	os.Setenv("APP_SECRET", "My Secret")
	id := "1234"
	s, err := utils.CreateToken(id)
	fmt.Printf("%v %v\n", s, err)

}

func TestParseToken(t *testing.T) {
	os.Setenv("APP_SECRET", "My Secret")
	id := "1234"
	s, err := utils.CreateToken(id)
	fmt.Printf("%v %v\n", s, err)

	claims := utils.GoClaims{}
	utils.ParseToken(s, &claims)
	// if status != http.StatusOK || merr != {
	// 	t.Errorf("Got: { status: %v, err: %v} Want: {status: %v, err: %v}", http)
	// }
	// fmt.Printf("%v %v\n", status, merr)
	// fmt.Printf("%+v\n", claims.Az)
	// fmt.Printf("%+v\n", claims.SClaims)
	// fmt.Printf("%+v\n", claims)
	if claims.SClaims.Subject != id {
		t.Errorf("got %v want %v", claims.SClaims.Subject, id)
	}
}
