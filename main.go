package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
)

var app = fiber.New()

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

var Users = make(map[int]User)

func init() {
	// check if users.txt exists and if not create it
	if _, err := os.Stat("users.txt"); os.IsNotExist(err) {
		file, err := os.Create("users.txt")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

func main() {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, " + c.Params("name"))
	})

	app.Post("/users/", func(c *fiber.Ctx) error {
		// get user from json in post request
		var user User
		if err := c.BodyParser(&user); err != nil {
			return err
		}
		user.ID = strconv.FormatInt(int64(len(Users))+1, 10)
		// add user to map
		Users[len(Users)] = user
		fmt.Println(Users)
		// write user to file
		file, err := os.OpenFile("users.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		file.WriteString(fmt.Sprintf("%s %s %s\n", user.ID, user.Name, user.Age))
		defer file.Close()
		// return user
		return c.JSON(user)
	})

	// get user from text file by ID
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		// get user from text file
		file, err := os.Open("users.txt")
		if err != nil {
			return err
		}
		defer file.Close()
		// read file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// split line by space
			line := strings.Split(scanner.Text(), " ")
			// create user
			user := User{
				ID:   line[0],
				Name: line[1],
				Age:  line[2],
			}
			// check if user ID matches
			if user.ID == id {
				return c.JSON(user)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
		return c.Status(404).SendString("User not found")
	})

	// get all users from fiel and return as json
	app.Get("/users", func(c *fiber.Ctx) error {
		// get user from text file
		file, err := os.Open("users.txt")
		if err != nil {
			return err
		}
		defer file.Close()
		// read file line by line
		scanner := bufio.NewScanner(file)
		var users []User
		for scanner.Scan() {
			// split line by space
			line := strings.Split(scanner.Text(), " ")
			// create user
			user := User{
				ID:   line[0],
				Name: line[1],
				Age:  line[2],
			}
			// add user to slice
			users = append(users, user)
		}
		return c.JSON(users)
	})

	app.Listen(":3000")
}
