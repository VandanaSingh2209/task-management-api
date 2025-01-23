package Middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type MiddlewareResponseBody struct {
	StatusCode int
	Message    string
	Body       struct{}
}

func RespondWithError(c *gin.Context, code int, message string) {
	response := MiddlewareResponseBody{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
	c.AbortWithStatusJSON(code, response)
}

var secretKey = []byte("your_secret_key")

func JWTTokenGenerate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iss": "AmplTest",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Error(userID, " @JWTTokenGenerate Error while signing the token:", err)
		return "", err
	}

	log.Info(userID, " @JWTTokenGenerate Token Signed succesfully:")
	return tokenString, nil
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			log.Error("@ValidateToken No token provided")
			//Controllers.NoDataFoundResponse(c, "No token provided")
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["sub"])
		}

		c.Next()
	}
}

var (
	clientRateLimiters = make(map[string]*rate.Limiter)
	mu                 sync.Mutex
)

func GetLimiter(clientIP string, rps float64, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exists := clientRateLimiters[clientIP]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rps), burst)
		clientRateLimiters[clientIP] = limiter
	}
	return limiter
}

func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := GetLimiter(clientIP, rps, burst)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
