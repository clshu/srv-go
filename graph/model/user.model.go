package model

import (
	"github.com/kamva/mgm/v3"
)

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	FirstName        string `json:"firstName" bson:"firstName"`
	LastName         string `json:"lastName" bson:"lastName"`
	Email            string `json:"email" bson:"email"`
	Password         string `json:"password" bson:"password"`
}
