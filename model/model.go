package model

import (
	
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Surname              string    `json:"surname"`
	Email                string    `json:"email"`
	BirthDate            time.Time `json:"birthDate"`
	Description          string    `json:"description"`
	ProfileImage         string    `json:"profileImage"`
	FriendRequestUserIDs []string  `json:"friendRequestUserIDs"`
	FriendIDs            []string  `json:"friendIds"`
	Password             string    `json:"password"`
	UserType             string    `json:"userType"`
	IsActivated          bool      `json:"isActivated"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
}

type UserDTO struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
	Password  string    `json:"password"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type UserCredentialsDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Description  string `json:"description"`
	ProfileImage string `json:"profileImage"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type ResetPasswordDTO struct {
	Password string `json:"password"`
}


type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserType string `json:"userType"`
	jwt.StandardClaims
}

type News struct{
	
	ID        string    `json:"id" bson:"id"`
	NewsID    string    `json:"news_id" bson:"news_id"`
	Content   string    `json:"content"  bson:"content"`
	Title 	  string    `json:"title" bson:"title"`
	Image  	  string    `json:"image" bson:"image"`
	CreatedAt time.Time `json:"created_at"  bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at"  bson:"updated_at"`
}