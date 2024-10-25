package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/liju-github/user-management/internal/mocks"
	"github.com/liju-github/user-management/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	app := fiber.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	userController := &UserController{userService: mockUserService}
	app.Post("/signup", userController.Signup)

	tests := []struct {
		name               string
		requestBody        models.UserSignupRequest
		expectedStatusCode int
		expectedResponse   string
		returnError        error
		expectSignupCall   bool
	}{
		{
			name: "successful signup",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "SecurePass@123",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusCreated,
			expectedResponse:   fmt.Sprintf(`{"message":"%v"}`, models.SignupSuccessful),
			returnError:        nil,
			expectSignupCall:   true,
		},
		{
			name: "user already exists",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "SecurePass@123",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusConflict,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.UserAlreadyExists),
			returnError:        errors.New(models.UserAlreadyExists),
			expectSignupCall:   true,
		},
		{
			name: "password too short",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "Abc1!",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, fmt.Sprintf(models.ErrPasswordLength, models.MinPasswordLength, models.MaxPasswordLength)),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "password too long",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "A1!a" + string(make([]byte, 70)),
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, fmt.Sprintf(models.ErrPasswordLength, models.MinPasswordLength, models.MaxPasswordLength)),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "password missing uppercase",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "password123!",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.ErrPasswordComplexity),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "password missing special character",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "Password123",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.ErrPasswordComplexity),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "invalid input - empty name",
			requestBody: models.UserSignupRequest{
				Name: "",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.ErrRequiredFieldsEmpty),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "invalid email format",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "invalid-email",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "SecurePass@123",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.ErrInvalidEmailFormat),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "invalid age",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         0,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "SecurePass@123",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.ErrNegativeAge),
			returnError:        nil,
			expectSignupCall:   false,
		},
		{
			name: "internal server error",
			requestBody: models.UserSignupRequest{
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Age:         30,
				Gender:      "Male",
				Address:     "123 Street",
				PhoneNumber: 1234567890,
				Password:    "SecurePass@123",
				ImageURL:    "http://image.url",
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedResponse:   `{"error":"some internal error"}`,
			returnError:        errors.New("some internal error"),
			expectSignupCall:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectSignupCall {
				if test.returnError != nil {
					mockUserService.EXPECT().Signup(&test.requestBody).Return(test.returnError).Times(1)
				} else {
					mockUserService.EXPECT().Signup(&test.requestBody).Return(nil).Times(1)
				}
			}

			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}

			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			assert.JSONEq(t, test.expectedResponse, buf.String())
		})
	}
}

func TestLogin(t *testing.T) {
	app := fiber.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	userController := NewUserController(mockUserService)
	app.Post("/login", userController.Login)

	tests := []struct {
		name               string
		requestBody        models.UserLoginRequest
		expectedStatusCode int
		mockError          error
		userBlocked        bool
		validateResponse   func(t *testing.T, response map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "SecurePass@123",
			},
			expectedStatusCode: fiber.StatusOK,
			mockError:          nil,
			userBlocked:        false,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.LoginSuccessful, response["message"])
				assert.NotEmpty(t, response["token"])
				assert.NotEmpty(t, response["refresh_token"])
				user, ok := response["user"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotNil(t, user)
				assert.Equal(t, false, user["is_blocked"])
			},
		},
		{
			name: "empty password",
			requestBody: models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			mockError:          nil,
			userBlocked:        false,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.InvalidInput, response["error"])
			},
		},
		{
			name: "empty email",
			requestBody: models.UserLoginRequest{
				Email:    "",
				Password: "SecurePass@123",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			mockError:          nil,
			userBlocked:        false,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.InvalidInput, response["error"])
			},
		},
		{
			name: "blocked user",
			requestBody: models.UserLoginRequest{
				Email:    "blocked@example.com",
				Password: "SecurePass@123",
			},
			expectedStatusCode: fiber.StatusUnauthorized,
			mockError:          nil,
			userBlocked:        true,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.UserIsBlocked, response["error"])
			},
		},
		{
			name: "wrong password",
			requestBody: models.UserLoginRequest{
				Email:    "test@example.com",
				Password: "WrongPass@123",
			},
			expectedStatusCode: fiber.StatusUnauthorized,
			mockError:          errors.New(models.InvalidInput),
			userBlocked:        false,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.InvalidInput, response["error"])
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.mockError != nil {
				mockUserService.EXPECT().
					Login(test.requestBody.Email, test.requestBody.Password).
					Return(nil, test.mockError)
			} else if test.requestBody.Email != "" && test.requestBody.Password != "" {
				user := &models.User{
					Email:     test.requestBody.Email,
					IsBlocked: test.userBlocked,
				}
				mockUserService.EXPECT().
					Login(test.requestBody.Email, test.requestBody.Password).
					Return(user, nil)
			}

			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}

			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			test.validateResponse(t, response)
		})
	}
}

func TestGetProfile(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		if userID != "" {
			c.Locals("ID", userID)
		}
		return c.Next()
	})

	app.Get("/profile", func(ctx *fiber.Ctx) error {
		ID, ok := ctx.Locals("ID").(string)
		if !ok || ID == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": models.InvalidID})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success", "ID": ID})
	})

	tests := []struct {
		name               string
		userID             string
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "Successful profile retrieval",
			userID:             "123",
			expectedStatusCode: fiber.StatusOK,
			expectedResponse:   map[string]interface{}{"message": "Success", "ID": "123"},
		},
		{
			name:               "Unauthorized access - missing user ID",
			userID:             "",
			expectedStatusCode: fiber.StatusUnauthorized,
			expectedResponse:   map[string]interface{}{"error": "Unauthorized or invalid user ID"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			req.Header.Set("Content-Type", "application/json")
			if test.userID != "" {
				req.Header.Set("X-User-ID", test.userID)
			}

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}
