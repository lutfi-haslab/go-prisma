//go:generate go run github.com/steebchen/prisma-client-go db push

package main

import (
	"api/db"
	_ "api/docs"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

var (
	client *db.PrismaClient
	ctx    context.Context
)

func init() {
	client = db.NewClient()
	client.Prisma.Connect()

	ctx = context.Background()
}

func disconnect() {
	client.Prisma.Disconnect()
}

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 0.0.0.0:8080
// @BasePath /
// @schemes http
func main() {
	defer disconnect()
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	app.Get("/", HealthCheck)
	app.Get("/create-post", CreatePost)
	app.Get("/get-post", GetPost)
	app.Get("/delete-post/:id", DeletePost)

	app.Listen("0.0.0.0:8080")
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": "Server is up and running",
	}, "application/problem+json")
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /create-post [get]
func CreatePost(c *fiber.Ctx) error {
	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set("Hi from Lutfii!"),
		db.Post.Published.Set(true),
		db.Post.Desc.Set(" is a database toolkit and makes databases easy."),
	).Exec(ctx)

	if err != nil {
		return nil
	}

	return c.JSON(createdPost)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /get-post [get]
func GetPost(c *fiber.Ctx) error {
	post, err := client.Post.FindMany().Exec(ctx)

	if err != nil {
		return nil
	}

	return c.JSON(post)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /delete-post/:id [get]
func DeletePost(c *fiber.Ctx) error {
	post, err := client.Post.FindUnique(db.Post.ID.Equals(c.Params("id"))).Delete().Exec(ctx)

	if err != nil {
		return nil
	}

	return c.JSON(post)
}

func Run() error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	// create a post
	createdPost, err := client.Post.CreateOne(
		db.Post.Title.Set("Hi from Prisma!"),
		db.Post.Published.Set(true),
		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
	).Exec(ctx)
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	fmt.Printf("created post: %s\n", result)

	// find a single post
	post, err := client.Post.FindUnique(
		db.Post.ID.Equals(createdPost.ID),
	).Exec(ctx)
	if err != nil {
		return err
	}

	result, _ = json.MarshalIndent(post, "", "  ")
	fmt.Printf("post: %s\n", result)

	// for optional/nullable values, you need to check the function and create two return values
	// `desc` is a string, and `ok` is a bool whether the record is null or not. If it's null,
	// `ok` is false, and `desc` will default to Go's default values; in this case an empty string (""). Otherwise,
	// `ok` is true and `desc` will be "my description".
	desc, ok := post.Desc()
	if !ok {
		return fmt.Errorf("post's description is null")
	}

	fmt.Printf("The posts's description is: %s\n", desc)

	return nil
}
