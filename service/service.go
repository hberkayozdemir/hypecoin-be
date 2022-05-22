package service

import (
	"github.com/hberkayozdemir/hypecoin-be/email"
	"github.com/hberkayozdemir/hypecoin-be/errors"
	"github.com/hberkayozdemir/hypecoin-be/model"
	"github.com/hberkayozdemir/hypecoin-be/repository"
	"github.com/hberkayozdemir/hypecoin-be/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service *Service) RegisterUser(userDTO model.UserDTO) (*model.User, error) {
	alreadyRegisteredUser, err := service.repository.GetUserByEmail(userDTO.Email)
	if alreadyRegisteredUser != nil {
		return nil, errors.UserAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:          utils.GenerateUUID(8),
		Name:        userDTO.Name,
		Surname:     userDTO.Surname,
		Email:       userDTO.Email,
		BirthDate:   userDTO.BirthDate,
		Password:    string(hashedPassword),
		UserType:    "user",
		IsActivated: false,
		CreatedAt:   time.Now().UTC().Round(time.Minute),
		UpdatedAt:   time.Now().UTC().Round(time.Minute),
		Latitude:    userDTO.Latitude,
		Longitude:   userDTO.Longitude,
	}

	newUser, err := service.repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	err = email.SendMail(newUser.Email, "Complete Registration", "Please click "+os.Getenv("PROD_HOSTNAME")+"/activation/"+user.ID)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) LoginUser(userCredentialsDTO model.UserCredentialsDTO) (*model.Token, *fiber.Cookie, error) {
	user, err := service.repository.GetUserByEmail(userCredentialsDTO.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.UserNotFound
	}

	if !user.IsActivated {
		return nil, nil, errors.Unauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentialsDTO.Password)); err != nil {
		return nil, nil, errors.WrongPassword
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaims{
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			Issuer: user.ID,
		},
	})

	token, err := claims.SignedString([]byte(""))

	if err != nil {
		return nil, nil, err
	}

	cookie := fiber.Cookie{
		Name:    "user-token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	return &model.Token{
		Token: token,
	}, &cookie, nil
}

func (service *Service) Activation(userID string) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.UserNotFound
	}

	if !user.IsActivated {
		user.IsActivated = true
	} else {
		return nil, errors.UserAlreadyActivated
	}

	return service.repository.UpdateUser(userID, *user)
}

func (service *Service) ForgotPassword(forgotPasswordDTO model.ForgotPasswordDTO) error {
	registeredUser, _ := service.repository.GetUserByEmail(forgotPasswordDTO.Email)
	if registeredUser == nil {
		return errors.UserNotFound
	}

	if !registeredUser.IsActivated {
		return errors.UserNotActivated
	}

	err := email.SendMail(forgotPasswordDTO.Email, "Reset Password", "You can reset your password click "+os.Getenv("REACT_HOSTNAME")+"/reset-password/"+registeredUser.ID+" here.")
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) ResetPassword(userID string, resetPasswordDTO model.ResetPasswordDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := service.repository.GetUser(userID)
	if err != nil {
		return errors.UserNotFound
	}

	user.Password = string(hashedPassword)

	_, err = service.repository.UpdateUser(userID, *user)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) GetUser(userID string) (*model.User, error) {
	return service.repository.GetUser(userID)
}

func (service *Service) GetNews() ([] *model.News, error) {

	return service.repository.GetNews()
}


func (service *Service) UpdateUser(userID string, updateUserDTO model.UpdateUserDTO) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}
	user.Description = updateUserDTO.Description
	user.ProfileImage = updateUserDTO.ProfileImage

	updatedUser, err := service.repository.UpdateUser(userID, *user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}


