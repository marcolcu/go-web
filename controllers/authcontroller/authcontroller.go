package authcontroller

import (
	"go-web-native/entities"
	"go-web-native/models/usermodel"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Session store
var store = session.New()

func Register(c *fiber.Ctx) error {
	var credentials Credentials
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if credentials.Name == "" || credentials.Email == "" || credentials.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please provide all fields"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	user := entities.User{
		Name:      credentials.Name,
		Email:     credentials.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := usermodel.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Automatically login the user
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: credentials.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Registration and login successful"})
}

func Login(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		// Render the login template
		return c.Render("views/auth/login.html", nil)
	}

	if c.Method() == fiber.MethodPost {
		var credentials Credentials
		if err := c.BodyParser(&credentials); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
		}

		if credentials.Email == "" || credentials.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please provide all fields"})
		}

		user, err := usermodel.GetUserByEmail(credentials.Email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
		}

		expirationTime := time.Now().Add(1 * time.Hour)
		claims := &Claims{
			Email: credentials.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		// Set cookie
		c.Cookie(&fiber.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		return c.Redirect("/home")
	}

	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"error": "Method not allowed"})
}

func Logout(c *fiber.Ctx) error {
	// Invalidate the token by setting an expired cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set an expiration in the past
		HTTPOnly: true,
		Path:     "/", // Ensure the cookie is removed across all paths
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.Redirect("/login?error=Unauthorized")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.Redirect("/login?error=Unauthorized")
		}
		return c.Status(fiber.StatusBadRequest).SendString("Bad request")
	}
	if !token.Valid {
		return c.Redirect("/login?error=Unauthorized")
	}

	return c.Next()
}
