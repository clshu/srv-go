package resolver

import "github.com/clshu/srv-go/graph/model"

func User2UserView(user *model.User) *model.UserView {
	return &model.UserView{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func CreateUserInput2User(data *model.CreateUserInput) *model.User {
	return &model.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}
}
