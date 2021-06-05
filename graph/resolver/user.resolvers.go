package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/clshu/srv-go/graph/model"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *mutationResolver) CreateUser(ctx context.Context, data model.CreateUserInput) (*model.UserView, error) {
	// index := mongo.IndexModel{Keys: bson.M{"email": 1}}
	// opt := options.Index(true)
	// opt.SetUnique()
	user := CreateUserInput2User(&data)
	err := mgm.Coll(user).Create(user)
	if err != nil {
		return nil, err
	}
	userView := User2UserView(user)
	return userView, nil
}

func (r *mutationResolver) LogIn(ctx context.Context, data model.LogInInput) (*model.AuthPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.UserView, error) {
	user := &model.User{}
	var results []*model.UserView
	// sortOptions := SortOptionsMap{"lastName": 1, "firstName": 1}
	// opts := CreateSortOptions(sortOptions)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "lastName", Value: 1}, {Key: "firstName", Value: 1}})

	mctx := mgm.Ctx()
	cur, err := mgm.Coll(user).Find(mctx, bson.M{}, findOptions)
	defer CursorClose(cur, mctx)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for cur.Next(mctx) {
		var user model.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		userView := User2UserView(&user)
		results = append(results, userView)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

func (r *queryResolver) Profile(ctx context.Context) (*model.UserView, error) {
	panic(fmt.Errorf("not implemented"))
}
