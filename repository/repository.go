package repository

import (
	"context"
	"log"
	"time"

	"github.com/hberkayozdemir/hypecoin-be/errors"
	"github.com/hberkayozdemir/hypecoin-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	MongoClient *mongo.Client
}

type UserEntity struct {
	ID                   string    `bson:"id"`
	Name                 string    `bson:"name"`
	Surname              string    `bson:"surname"`
	Email                string    `bson:"email"`
	BirthDate            time.Time `bson:"birthDate"`
	Description          string    `bson:"description"`
	ProfileImage         string    `bson:"profileImage"`
	FriendRequestUserIDs []string  `bson:"friendRequestUserIDs"`
	FriendIDs            []string  `bson:"friendIds"`
	Password             string    `bson:"password"`
	UserType             string    `bson:"userType"`
	IsActivated          bool      `bson:"isActivated"`
	CreatedAt            time.Time `bson:"createdAt"`
	UpdatedAt            time.Time `bson:"updatedAt"`
	Latitude             float64   `bson:"latitude"`
	Longitude            float64   `bson:"longitude"`
}

type PostEntity struct {
	ID              string    `bson:"id"`
	UserID          string    `bson:"userId"`
	Description     string    `bson:"description"`
	Image           string    `bson:"image"`
	IsPrivate       bool      `bson:"isPrivate"`
	WhoLikesUserIDs []string  `bson:"whoLikesUserIds"`
	CommentIDs      []string  `bson:"commentIds"`
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}

type CommentEntity struct {
	ID        string    `bson:"id"`
	UserID    string    `bson:"userId"`
	PostID    string    `bson:"postId"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func NewRepository(uri string) *Repository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	defer cancel()
	client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return &Repository{client}
}

func (repository *Repository) RegisterUser(user model.User) (*model.User, error) {
	collection := repository.MongoClient.Database("hypecoin").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userEntity := convertUserModelToUserEntity(user)

	_, err := collection.InsertOne(ctx, userEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetUser(userEntity.ID)
}

func (repository *Repository) GetUser(userID string) (*model.User, error) {
	collection := repository.MongoClient.Database("hypecoin").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": userID}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.UserNotFound
	}

	userEntity := UserEntity{}
	err := cur.Decode(&userEntity)

	if err != nil {
		return nil, err
	}

	user := convertUserEntityToUserModel(userEntity)

	return &user, nil
}

func (repository *Repository) GetUserByEmail(email string) (*model.User, error) {
	collection := repository.MongoClient.Database("hypecoin").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, errors.UserNotFound
	}

	if cur == nil {
		return nil, errors.UserNotFound
	}

	userEntity := UserEntity{}
	err := cur.Decode(&userEntity)

	if err != nil {
		return nil, err
	}

	user := convertUserEntityToUserModel(userEntity)

	return &user, nil
}

func (repository *Repository) UpdateUser(userID string, user model.User) (*model.User, error) {
	collection := repository.MongoClient.Database("hypecoin").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": userID}

	userEntity := convertUserModelToUserEntity(user)

	cur := collection.FindOneAndReplace(ctx, filter, userEntity)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.UserNotFound
	}

	return repository.GetUser(userID)

}


func (repository *Repository) GetNews() ([]*model.News, error) {
	collection := repository.MongoClient.Database("hypecoin").Collection("news")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur,_:= collection.Find(ctx,nil)
	if cur.Err() != nil {
		return nil, cur.Err()
	}

	var news []*model.News
	for cur.Next(ctx) {
		newEntity :=&model.News{}
		err := cur.Decode(newEntity)
		if err != nil {
			return nil, err
		}
		news = append(news,newEntity)
	if cur == nil {
		return nil, errors.NewsNotFound
	}
	}
	return news, nil
}





func (repository *Repository) GetUsersByIDList(userIDs []string) ([]model.User, error) {
	collection := repository.MongoClient.Database("hypecoin").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M
	if len(userIDs) == 0 {
		filter = bson.M{"id": bson.M{"$in": []string{}}}
	} else {
		filter = bson.M{"id": bson.M{"$in": userIDs}}
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for cur.Next(ctx) {
		userEntity := UserEntity{}
		err := cur.Decode(&userEntity)
		if err != nil {
			return nil, err
		}
		users = append(users, convertUserEntityToUserModel(userEntity))
	}
	return users, nil
}

func convertUserModelToUserEntity(user model.User) UserEntity {
	return UserEntity{
		ID:                   user.ID,
		Name:                 user.Name,
		Surname:              user.Surname,
		Email:                user.Email,
		BirthDate:            user.BirthDate,
		Description:          user.Description,
		ProfileImage:         user.ProfileImage,
		FriendRequestUserIDs: user.FriendRequestUserIDs,
		FriendIDs:            user.FriendIDs,
		Password:             user.Password,
		UserType:             user.UserType,
		IsActivated:          user.IsActivated,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
		Latitude:             user.Latitude,
		Longitude:            user.Longitude,
	}
}

func convertUserEntityToUserModel(userEntity UserEntity) model.User {
	return model.User{
		ID:                   userEntity.ID,
		Name:                 userEntity.Name,
		Surname:              userEntity.Surname,
		Email:                userEntity.Email,
		BirthDate:            userEntity.BirthDate,
		Description:          userEntity.Description,
		ProfileImage:         userEntity.ProfileImage,
		FriendRequestUserIDs: userEntity.FriendRequestUserIDs,
		FriendIDs:            userEntity.FriendIDs,
		Password:             userEntity.Password,
		UserType:             userEntity.UserType,
		IsActivated:          userEntity.IsActivated,
		CreatedAt:            userEntity.CreatedAt,
		UpdatedAt:            userEntity.UpdatedAt,
		Latitude:             userEntity.Latitude,
		Longitude:            userEntity.Longitude,
	}
}



