package utils

import (
	"time"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

// TempPassword holds temperary password and expiration date
type TempPassword struct {
	Str     string `json:"str" bson:"str"`
	ExpDate int64  `json:"expDate" bson:"expDate"`
}

// PasswordMinLenth is password minmun length
const passwordMinLenth = 8

// TempPasswordExpLength is Temp password expiration length
const tempPasswordExpLength = 24 * 60 * 60 * 1000 // 24 hours in milliseconds

// CreateHashedPassword creates a hashed password
func CreateHashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// ComparePassword compares hashed and plain text password
func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CreateTempPassword creates temperary password
func CreateTempPassword() (TempPassword, error) {
	// 10 characters long
	password, err := password.Generate(10, 1, 1, false, false)
	if err != nil {
		return TempPassword{}, err
	}
	// Javascript in browser code represents UTC time in milliseconds
	// Convert Unix() seconds to milliseconds
	expDate := time.Now().Unix()*1000 + tempPasswordExpLength

	ret := TempPassword{
		Str:     password,
		ExpDate: expDate,
	}
	return ret, nil
}

// CompareTempPassword compares temperary passwords
func CompareTempPassword(password string, tempPassword TempPassword) bool {
	return password == tempPassword.Str
}
