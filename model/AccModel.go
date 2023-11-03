package model

import (
	"GoAPI/initializer"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

type AuthInput struct {
	UserName string `json:"username" bson:"username" form:"username"`
	Password string `json:"password" bson:"password" form:"password"`
	Email    string `json:"email" bson:"email" form:"email"`
}

type Account struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email" `
}

// hàm thêm yêu cầu đăng nhập từ form html
func (account *Account) Save() (*Account, error) {
	collection := initializer.UserDB()
	defer initializer.UserDB()
	_, err := collection.InsertOne(context.TODO(), account)
	if err != nil {
		return &Account{}, err
	}
	return account, nil
}

func IsAccountValid(username, password string) bool {
	return false
}

// hash mật khẩu
func (account *Account) PreProcess() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(passwordHash)
	account.UserName = html.EscapeString(strings.TrimSpace(account.UserName))
	return nil
}

// thêm method Valid cho User
func (account *Account) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
}

func FindAccount(target string) (Account, error) {
	collection := initializer.UserDB()
	defer initializer.UserDB()
	var user Account
	filter := bson.M{"$or": []bson.M{
		{"username": target},
		{"email": target},
	},
	}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return Account{}, err
	}
	return user, nil
}
