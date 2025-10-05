package auth

import (
	"go-auth/internal/database"
	"go-auth/internal/response"
	"go-auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrUserExists      = "User already exists"
	ErrUserNotFound    = "User does not exist"
	ErrInvalidPassword = "Invalid password"
	ErrTokenFailure    = "Failed to create token"
)

func SignUp(db *database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			response.Respond(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Respond(c, http.StatusInternalServerError, ErrTokenFailure, nil)
			return
		}

		user := &models.User{
			Email:    body.Email,
			Password: string(hashedPassword),
			Role:     "visitor",
		}

		if err := db.CreateUser(user); err != nil {
			response.Respond(c, http.StatusConflict, ErrUserExists, nil)
			return
		}

		response.Respond(c, http.StatusCreated, "User created successfully", nil)
	}
}

func Login(db *database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			response.Respond(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		user, err := db.FindUserByEmail(body.Email)
		if err != nil {
			response.Respond(c, http.StatusInternalServerError, "Something went wrong", nil)
			return
		}
		if user.ID == 0 {
			response.Respond(c, http.StatusNotFound, ErrUserNotFound, nil)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			response.Respond(c, http.StatusUnauthorized, ErrInvalidPassword, nil)
			return
		}

		accessToken, err := GenerateAccessTokenString(user)
		if err != nil {
			response.Respond(c, http.StatusBadRequest, ErrTokenFailure, nil)
			return
		}

		refreshToken, err := GenerateAccessTokenString(user)
		if err != nil {
			response.Respond(c, http.StatusBadRequest, ErrTokenFailure, nil)
			return
		}

		response.Respond(c, http.StatusOK, "Login successful",
			gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"expires_in":    900, // 15mins
			})
	}
}

func RefreshAccessToken(db *database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			response.Respond(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		claims, err := ParseToken(body.RefreshToken)
		if err != nil {
			response.Respond(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		userId := claims.Subject
		user, err := db.FindUserById(userId)
		if err != nil {
			response.Respond(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		accessToken, err := GenerateAccessTokenString(user)
		if err != nil {
			response.Respond(c, http.StatusInternalServerError, ErrTokenFailure, nil)
		}
		response.Respond(c, http.StatusOK, "Token Refreshed",
			gin.H{
				"access_token": accessToken,
				"expires_in":   900,
			})
	}
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	response.Respond(c, http.StatusOK, "Logged out successfully", nil)
}

func Validate(c *gin.Context) {
	response.Respond(c, http.StatusOK, "Logged in", nil)
}
