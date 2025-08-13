package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, email, username string) (string, error) {
	tokenTTLStr := os.Getenv("TOKEN_TTL")
	tokenTTL, err := strconv.Atoi(tokenTTLStr)
	if err != nil {
		tokenTTL = 2000 // default 2000 minutes
	}

	claims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenTTL) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_PRIVATE_KEY")
	return token.SignedString([]byte(secretKey))
}

func ValidateJWT(c *gin.Context) error {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return errors.New("authorization token required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		secretKey := os.Getenv("JWT_PRIVATE_KEY")
		return []byte(secretKey), nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}

	c.Set("user_id", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("username", claims.Username)

	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	
	bearerToken := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}
	
	return ""
}

func ExtractUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	
	uid, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID format")
	}
	
	return uid, nil
}