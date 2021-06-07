package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

const domain string = "thermovisionusa.com"
const algorithm string = "HS256"

// const Domain string = strings.ToUpper(domain)
const oneDay int64 = 24 * 60 * 60 // in seconds

type Authz struct {
	Rol []string `json:"rol"` // roles
	Per []string `json:"per"` // permmisions
}

type GoClaims struct {
	Az      Authz              `json:"az"`
	SClaims jwt.StandardClaims `json:"sclaims"`
}

type TokenId struct {
	ID string `json:"id"`
	Az Authz  `json:"az"`
}

// Valid is to implement Claim interface
func (c *GoClaims) Valid() error {
	return nil
}

func CreateToken(id string) (string, error) {
	secret := os.Getenv("APP_SECRET")
	if secret == "" {
		err := fmt.Errorf("Do Not Find Secret for Signing Token")
		return "", err
	}
	az := Authz{}
	now := time.Now().Unix()
	mapClaims := GoClaims{
		Az: az,
		SClaims: jwt.StandardClaims{
			// Id:        "0",
			Subject:   id,
			IssuedAt:  now,
			ExpiresAt: now + oneDay,
			Issuer:    domain,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &mapClaims)
	// fmt.Printf("%v\n", token)
	ss, err := token.SignedString([]byte(secret))

	return ss, err
}

func ParseToken(tokenString string) (*TokenId, error) {
	// mClaims := GoClaims{}
	secret := os.Getenv("APP_SECRET")
	// secret = "5678"
	if secret == "" {
		err := fmt.Errorf("Do Not Find Secret for Signing Token")
		return nil, err
	}
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &GoClaims{}, func(token *jwt.Token) (interface{}, error) {

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}

	// Check against alg=none hack
	if token.Method.Alg() != algorithm {
		return nil, fmt.Errorf("Algorithm is incorrect")
	}
	claims, ok := token.Claims.(*GoClaims)
	if !ok {
		return nil, fmt.Errorf("Token is corrupted")
	}
	tokenID := TokenId{ID: claims.SClaims.Subject, Az: claims.Az}

	return &tokenID, nil
}

// func GetTokenId(req *http.Request) (*TokenId, int, error) {
// 	str := req.Header.Get("Authorization")
// 	if str == "" {
// 		return nil, http.StatusUnauthorized, fmt.Errorf("No Token in Header")
// 	}

// 	bearer := strings.Split(str, " ")

// 	if strings.ToLower(bearer[0]) != "bearer" {
// 		return nil, http.StatusUnauthorized, fmt.Errorf("No Bearer token")
// 	}
// 	claims := GoClaims{}

// 	status, err := ParseToken(bearer[1], &claims)

// 	if err != nil {
// 		return nil, status, err
// 	}
// 	tokenID := TokenId{ID: claims.SClaims.Subject, Az: claims.Az}
// 	return &tokenID, http.StatusOK, nil

// }
