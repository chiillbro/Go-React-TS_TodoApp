package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Todo is a struct that represents a todo item

type Todo struct {
		ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Completed bool `json: "completed"`
		Body string `json: "body"`
	}


	var collection *mongo.Collection


func main() {
	fmt.Println("Hello world")

	// loading .env file

	if os.Getenv("ENV") != "production" {

		err := godotenv.Load(".env")
		if err !=  nil {
			log.Fatal("Error loading .env file")
		}

	}
	

	// Getting the PORT from .env file
	PORT := os.Getenv("PORT")


	
	// Getting the MONGO_URI from .env file
	MONGO_URI :=os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err =  client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("Go_Todo_DB").Collection("todos")


	app := fiber.New()

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}


	// Cors for development

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:5174",
	// 	AllowHeaders: "origin, content-type, accept",
	// }))


	// Get all todos

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)


	log.Fatal(app.Listen(":" + PORT))
	
}



// Get all todos function
func getTodos(c *fiber.Ctx) error {


	// Initializing empty todos array
   var todos []Todo
	

	cursor, err := collection.Find(context.Background(),bson.M{})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting todos"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error getting todos"})
	}
	
	todos = append(todos, todo)

	}
	return c.JSON(todos)
}


// Create a new todo function

func createTodo(c *fiber.Ctx)  error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating todo"})
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)	
}



	// Update a todo

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set" : bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err !=  nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating todo"})
	}
	var updatedTodo Todo

	err = collection.FindOne(context.Background(), filter).Decode(&updatedTodo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error Fetching updated todo"})
	}
	
	return c.Status(200).JSON(updatedTodo)
}


// Delete a todo

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err :=primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo id"})
	}

	filter := bson.M{"_id": objectId}


	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting todo"})
	}

	var todos []Todo

	cursor, err := collection.Find(context.Background(),bson.M{})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting todos"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error getting todos"})
	}
	
	todos = append(todos, todo)

	}


	return c.Status(200).JSON(fiber.Map{"todos": todos ,"msg": "Todo deleted"})

}




























// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// )

// // Todo is a struct that represents a todo item

// type Todo struct {
// 	ID int `json: "id"`
// 	Completed bool `json: "completed"`
// 	Body string `json: "body"`
// }

// func main() {
// 	fmt.Println("how are you")
// 	app := fiber.New()

// 	// loading .env file
// 	err := godotenv.Load(".env")

// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

//     // Set the port

// 	PORT := os.Getenv("PORT")

// 	// Initializing todos array
// 	todos := []Todo{}

// 	// Get all todos
// 	app.Get("/api/todos", func(c *fiber.Ctx) error {
// 		// return c.Status(200).JSON(fiber.Map{"msg": "Hello World"})
// 		return c.Status(200).JSON(todos)
//    })

//    // Create a new todo
// 	app.Post("/api/todos", func(c *fiber.Ctx) error {
// 	todo := &Todo{}

// 	if err := c.BodyParser(todo); err != nil {
// 		return err
// 	}

// 	if todo.Body == "" {
// 		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
// 	}

// 	todo.ID = len(todos) + 1
// 	todos = append(todos, *todo)

// 	return c.Status(201).JSON(todo)
//    })

//    // Update a todo

//    app.Patch("/api/todos/:id", func (c *fiber.Ctx) error  {

// 	id := c.Params("id")

// 	for i, todo := range todos {
// 		if fmt.Sprint(todo.ID) == id {
// 			todos[i].Completed = true

// 			return c.Status(200).JSON(todos[i])
// 		}
// 	}
// 	return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
//    })

//    // Delete a todo

//    app.Delete("/api/todos/:id", func (c *fiber.Ctx) error  {

// 	id :=  c.Params("id")
// 	for i, todo :=  range todos {
// 		if fmt.Sprint(todo.ID) == id {
// 			todos = append(todos[:i], todos[i+1:]...)
// 			return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted"})
// 	}
//    }

//    return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
//    })

// 	log.Fatal(app.Listen(":" + PORT))
// }