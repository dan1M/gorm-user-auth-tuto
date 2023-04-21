package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dan1M/gorm-user-auth-tuto/config"
	"github.com/dan1M/gorm-user-auth-tuto/model"
	"github.com/dan1M/gorm-user-auth-tuto/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	RtService   *service.RtService
	UserService *service.UserService
	*config.Config
}

func NewAuthHandler(rtService *service.RtService, userService *service.UserService, config *config.Config) *AuthHandler {
	return &AuthHandler{rtService, userService, config}
}

func (auth *AuthHandler) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(auth.JWT_SECRET))
}

func (auth *AuthHandler) Login(c *gin.Context) {
	var loginDTO *model.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.UserService.GetByEmail(loginDTO.Email)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = user.CheckPassword(loginDTO.Password)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jwt, err := auth.GenerateToken(user)
	if err != nil {
		fmt.Println(err)
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(400, gin.H{"error": "Invalid password"})
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		return
	}

	rt, err := auth.RtService.CreateRT(&model.RtCreateDTO{
		UserID: int(user.ID),
		Ip:     c.ClientIP(),
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", jwt, 3600, "/", "*", true, true)
	c.SetCookie("rt", rt.Hash, 3600, "/", "*", true, true)

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   jwt,
		"rt":      rt.Hash,
		"user":    user,
	})
}

func (auth *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		jwtCookie, err := c.Cookie("jwt")
		if err == http.ErrNoCookie {
			authHeader := c.GetHeader("Authorization")
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) != 2 {
				c.JSON(401, gin.H{"error": "Invalid authorization header"})
				c.Abort()
				return
			}
			jwtCookie = splitToken[1]

			if jwtCookie == "" {
				c.JSON(401, gin.H{"error": "No token provided"})
				c.Abort()
				return
			}
		}
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			}
			return []byte(auth.JWT_SECRET), nil
		})
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userId := token.Claims.(jwt.MapClaims)["id"].(float64)
		user, err := auth.UserService.GetUser(int(userId))
		if err != nil {
			c.JSON(401, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()

		// after request
	}
}
