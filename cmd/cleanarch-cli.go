package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type ProjectConfig struct {
	Name   string
	Module string
}

func main() {
	config := ProjectConfig{}
	flag.StringVar(&config.Name, "name", "", "Name of the project")
	flag.StringVar(&config.Module, "module", "", "Go module name (e.g., github.com/username/project)")
	flag.Parse()

	if config.Name == "" || config.Module == "" {
		fmt.Println("Please provide both project name and module name")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := generateProject(config)
	if err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project generated successfully!")
}

func generateProject(config ProjectConfig) error {
	// Create project directory
	err := os.MkdirAll(config.Name, 0755)
	if err != nil {
		return err
	}

	// Change to project directory
	err = os.Chdir(config.Name)
	if err != nil {
		return err
	}

	// Create directory structure
	dirs := []string{
		"cmd/api",
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler",
		"pkg/config",
		"pkg/database",
		"pkg/auth",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// Generate files
	files := map[string]string{
		"cmd/api/main.go":                        mainTemplate,
		"internal/domain/user.go":                userDomainTemplate,
		"internal/usecase/user_usecase.go":       userUsecaseTemplate,
		"internal/repository/user_repository.go": userRepositoryTemplate,
		"internal/handler/user_handler.go":       userHandlerTemplate,
		"pkg/config/config.go":                   configTemplate,
		"pkg/database/mongodb.go":                mongodbTemplate,
		"pkg/auth/jwt.go":                        jwtTemplate,
		"go.mod":                                 goModTemplate,
		".env":                                   envTemplate,
	}

	for file, tmpl := range files {
		err := generateFile(file, tmpl, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateFile(filename, tmplContent string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New(filepath.Base(filename)).Parse(tmplContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, data)
}

// Templates for each file (add these at the end of the file)
const mainTemplate = `package main

import (
	"log"
	"{{.Module}}/internal/handler"
	"{{.Module}}/internal/repository"
	"{{.Module}}/internal/usecase"
	"{{.Module}}/pkg/auth"
	"{{.Module}}/pkg/config"
	"{{.Module}}/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up JWT secret
	auth.SetJWTSecret(cfg.JWTSecret)

	// Connect to MongoDB
	db := database.ConnectMongoDB(cfg.MongoURI)

	// Set up repositories
	userRepo := repository.NewUserRepository(db)

	// Set up use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Set up handlers
	userHandler := handler.NewUserHandler(userUseCase)

	// Set up Fiber app
	app := fiber.New()

	// Set up routes
	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)

	// Start server
	log.Fatal(app.Listen(":8080"))
}
`

const userDomainTemplate = `package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID ` + "`json:\"id\" bson:\"_id,omitempty\"`" + `
	Username string             ` + "`json:\"username\" bson:\"username\"`" + `
	Email    string             ` + "`json:\"email\" bson:\"email\"`" + `
	Password string             ` + "`json:\"-\" bson:\"password\"`" + `
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}

type UserUseCase interface {
	Register(user *User) error
	Login(email, password string) (string, error)
}
`

const userUsecaseTemplate = `package usecase

import (
	"errors"
	"{{.Module}}/internal/domain"
	"{{.Module}}/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: repo}
}

func (uc *userUseCase) Register(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return uc.userRepo.Create(user)
}

func (uc *userUseCase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}
`

const userRepositoryTemplate = `package repository

import (
	"context"
	"{{.Module}}/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) Create(user *domain.User) error {
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
`

const userHandlerTemplate = `package handler

import (
	"{{.Module}}/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(useCase domain.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: useCase}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := h.userUseCase.Register(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Registration failed"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string ` + "`json:\"email\"`" + `
		Password string ` + "`json:\"password\"`" + `
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := h.userUseCase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"token": token})
}
`

const configTemplate = `package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	JWTSecret string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
`

const mongodbTemplate = `package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("myapp")
}
`

const jwtTemplate = `package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString(jwtSecret)
}
`

const goModTemplate = `module {{.Module}}

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofiber/fiber/v2 v2.38.1
	github.com/joho/godotenv v1.4.0
	go.mongodb.org/mongo-driver v1.10.3
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
)
`

const envTemplate = `MONGO_URI=your_mongodb_atlas_connection_string
JWT_SECRET=your_jwt_secret
`
