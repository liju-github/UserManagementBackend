package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Signup(user *models.UserSignupRequest) error {
	existingUser, _ := s.userRepo.FindUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New(models.UserAlreadyExists)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// verificationToken := generateToken()

	newUser := &models.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   true,
		Age:          user.Age,
		Gender:       user.Gender,
		PhoneNumber:  user.PhoneNumber,
		Address:      user.Address,
		// VerificationToken:  verificationToken,
		// VerificationExpiry: time.Now().Add(24 * time.Hour).Unix(),
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err
	}

	// s.sendVerificationEmail(newUser.Email, verificationToken)

	return nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New(models.UserDoesntExist)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified")
	}

	return user, nil
}

func (s *UserService) Logout(userID string) error {
	return nil
}

func (s *UserService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindUserByVerificationToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	if time.Now().Unix() > user.VerificationExpiry {
		return errors.New("token expired")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	user.VerificationExpiry = 0

	return s.userRepo.UpdateUser(user)
}

func (s *UserService) ResendVerification(email string) error {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return errors.New(models.UserDoesntExist)
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	verificationToken := generateToken()
	user.VerificationToken = verificationToken
	user.VerificationExpiry = time.Now().Add(24 * time.Hour).Unix()

	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return s.sendVerificationEmail(user.Email, verificationToken)
}

func (s *UserService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return errors.New(models.UserDoesntExist)
	}

	resetToken := generateToken()
	resetExpiry := time.Now().Add(1 * time.Hour).Unix()

	passwordReset := &models.PasswordReset{
		UserID:     user.ID,
		ResetToken: resetToken,
		Expiry:     resetExpiry,
	}

	if err := s.userRepo.CreatePasswordReset(passwordReset); err != nil {
		return err
	}

	return s.sendPasswordResetEmail(user.Email, resetToken)
}

func (s *UserService) ConfirmPasswordReset(token, newPassword string) error {
	passwordReset, err := s.userRepo.FindPasswordResetByToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	if time.Now().Unix() > passwordReset.Expiry {
		return errors.New("token expired")
	}

	user, err := s.userRepo.FindUserByID(passwordReset.UserID)
	if err != nil {
		return errors.New(models.UserDoesntExist)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)

	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return s.userRepo.DeletePasswordReset(passwordReset.ID)
}

func (s *UserService) GetProfile(userID string) (*models.User, error) {
	return s.userRepo.FindUserByID(userID)
}

func (s *UserService) UpdateProfile(userID string, email string, req *models.UserUpdateRequest) error {
	// Find the existing user by ID
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return errors.New(models.UserDoesntExist)
	}

	// Update only the mutable fields
	user.Name = req.Name
	user.Age = req.Age
	user.Gender = req.Gender
	user.Address = req.Address

	fmt.Println("to be updated ", user)

	// Call the repository to update the user in the database
	return s.userRepo.UpdateUser(user) // Pass the updated user object
}

func (s *UserService) UploadProfilePicture(userID, cdnURL string) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return errors.New(models.UserDoesntExist)
	}

	user.ImageURL = cdnURL
	return s.userRepo.UpdateUser(user)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *UserService) sendVerificationEmail(email, token string) error {
	from := "lijuthomasliju03@gmail.com"
	password := "ciwg zzwn gpbs dekx"
	to := []string{"h4ze07@gmail.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: Email Verification\r\n" +
		"\r\n" +
		"Please verify your email by clicking the following link:\r\n" +
		"http://example.com/verify?token=" + token)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) sendPasswordResetEmail(email, token string) error {
	return s.sendVerificationEmail(email, token)
}
