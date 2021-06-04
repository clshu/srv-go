package model

import (
	"strings"

	"github.com/clshu/srv-go/utils"
)

// Creating : an preop to create hashed password and lower case email
func (u *User) Creating() error {
	// Call the DefaultModel Creating hook
	if err := u.DefaultModel.Creating(); err != nil {
		return err
	}

	if u.Email != "" {
		u.Email = strings.ToLower(u.Email)
	}
	if u.Password != "" {
		hash, err := utils.CreateHashedPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = string(hash)
	}

	return nil
}
