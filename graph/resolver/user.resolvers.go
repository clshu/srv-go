package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/clshu/srv-go/graph/generated"
	"github.com/clshu/srv-go/graph/model"
	"github.com/clshu/srv-go/middleware/auth"
	"github.com/clshu/srv-go/utils"
	mgm "github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	user := &model.User{}
	mctx := mgm.Ctx()
	// Normalize the email to lower case
	email := strings.ToLower(data.Email)
	password := data.Password

	result := mgm.Coll(user).FindOne(mctx, bson.M{"email": email})
	if result.Err() != nil {
		if strings.Contains(result.Err().Error(), "no documents in result") {
			// log.Printf("%s - %s", result.Err(), email)
			return nil, fmt.Errorf("Unable to login")
		}
		return nil, result.Err()
	}

	err := result.Decode(user)
	if err != nil {
		// log.Print(err)
		return nil, fmt.Errorf("Unable to login")
	}
	if user.Email != email {
		// log.Print(err)
		return nil, fmt.Errorf("Unable to login")
	}

	// Compare hashed password to plain text password
	err = utils.ComparePassword(user.Password, password)
	if err != nil {
		// log.Print(err)
		return nil, fmt.Errorf("Unable to login")
	}
	var token string
	token, err = utils.CreateToken(user.ID.Hex())
	if err != nil {
		// log.Print(err)
		return nil, err
	}
	ret := &model.AuthPayload{
		User:  User2UserView(user),
		Token: token,
	}
	return ret, nil
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
		// log.Print(err)
		return nil, err
	}
	for cur.Next(mctx) {
		var user model.User
		err := cur.Decode(&user)
		if err != nil {
			// log.Print(err)
			return nil, err
		}
		userView := User2UserView(&user)
		results = append(results, userView)
	}

	if err := cur.Err(); err != nil {
		// log.Print(err)
		return nil, err
	}

	return results, nil
}

func (r *queryResolver) Profile(ctx context.Context) (*model.UserView, error) {
	tokenId, err := auth.GetTokenIdFromCtx(ctx)
	if err != nil {
		// Header or token parsing errors
		return nil, err
	}
	// log.Println(tokenId)
	if tokenId == nil {
		// No Authorization token
		return nil, fmt.Errorf(("Unauthorized"))
	}
	user := &model.User{}
	id, _ := primitive.ObjectIDFromHex(tokenId.ID)
	mctx := mgm.Ctx()
	result := mgm.Coll(user).FindOne(mctx, bson.M{"_id": id})
	if result.Err() != nil {
		if strings.Contains(result.Err().Error(), "no documents in result") {
			// log.Printf("%s - %s", result.Err(), email)
			return nil, fmt.Errorf("Profile Not Found")
		}
		return nil, result.Err()
	}

	err = result.Decode(user)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return User2UserView(user), nil
}

func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
