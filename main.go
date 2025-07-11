package main

import (
	"database/sql" // Added for database interaction
	"fmt"          // Added for printing messages
	"log"          // Added for logging errors
	"net/http"
	"time" // Added for setting connection pool parameters

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// Define a struct to represent the data we expect for a new bike
type Bike struct {
	Name      string `json:"name" binding:"required"`
	WheelSize int    `json:"wheel_size" binding:"required,gte=10"`
	Color     string `json:"color,omitempty"`
}

// Global variable for the database connection pool (we'll discuss better ways like dependency injection in Week 7)
var db *sql.DB

func main() {
	// --- Database Connection Logic START ---
	var err error // Declare err here to be used for db connection and router
	// Replace with your actual connection details
	connStr := "postgres://postgres:ParleG%40123@localhost:5432/postgres?sslmode=disable" //
	// It's good practice to use a specific database for your project, e.g., "finspeed_db"
	// You might need to create it first using psql: CREATE DATABASE finspeed_db;
	// Then change "postgres" in connStr to "finspeed_db"

	// Open a connection pool
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to open database connection: %v\n", err)
	}
	// It's good practice to set connection pool parameters
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")
	// --- Database Connection Logic END ---

	// 1. Create a new Gin router with default middleware (logger, recovery)
	router := gin.Default()

	// 2. Define a route for GET requests to the "/ping" path
	router.GET("/ping", func(c *gin.Context) {
		// This function is the "handler" for the "/ping" route

		// 3. Send a JSON response with status 200 OK
		c.JSON(http.StatusOK, gin.H{ // gin.H is a shortcut for map[string]interface{}
			"message": "pong",
		})
	})

	// Add this new route below the /ping route
	// Route for getting a specific bike by ID
	router.GET("/api/bikes/:id", func(c *gin.Context) {
		// 1. Get the value of the "id" parameter from the URL path
		bikeID := c.Param("id") // "id" matches the :id in the route definition

		// 2. Send a JSON response including the extracted ID
		c.JSON(http.StatusOK, gin.H{
			"message":  "Fetching details for bike ID",
			"bike_id":  bikeID, // Include the captured ID in the response
		})
	})

	// Add this new route below the /api/bikes/:id route
		// Route for getting a list of bikes, potentially filtered by query parameters
		router.GET("/api/bikes", func(c *gin.Context) {
			// 1. Get query parameters: "type" and "color"
			// c.Query("key") returns the value, or an empty string "" if not present.
			// c.DefaultQuery("key", "defaultValue") returns the value, or "defaultValue" if not present.

			bikeType := c.DefaultQuery("type", "any") // Get 'type', default to 'any' if not provided
			bikeColor := c.Query("color")           // Get 'color', will be "" if not provided

			// Prepare the response message
			response := gin.H{
				"message":      "Fetching list of bikes",
				"filter_type":  bikeType,
				"filter_color": "not specified", // Default color message
			}

			// Only add color to response if it was actually provided
			if bikeColor != "" {
				response["filter_color"] = bikeColor
			}

			// 2. Send the JSON response
			c.JSON(http.StatusOK, response)
		})

	// ... (your GET routes) ...

		// Route for creating a new bike (handles POST requests)
		router.POST("/api/bikes", func(c *gin.Context) {
			var newBike Bike // Create a variable of our Bike struct type

			// 1. Bind the incoming JSON from the request body to the newBike struct
			// If there's an error (e.g., malformed JSON, missing required fields),
			// ShouldBindJSON will populate the error.
			if err := c.ShouldBindJSON(&newBike); err != nil {
				// If binding fails, send a 400 Bad Request response with the error details
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return // Important to return after sending an error response
			}

			// 2. If binding is successful, process the data (for now, just send it back)
			// In a real application, you would save newBike to a database here.
			c.JSON(http.StatusCreated, gin.H{ // 201 Created is a more appropriate status for successful creation
				"message":    "Bike created successfully!",
				"bike_name":  newBike.Name,
				"wheel_size": newBike.WheelSize,
				"color":      newBike.Color, // Will be empty if not provided
			})
		})

		// router.Run(":8080") ...
	
	// 4. Start the HTTP server and listen for requests on port 8080
	// By default, it listens on localhost (127.0.0.1)
	err = router.Run(":8080") // Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		// Handle error if the server fails to start (e.g., port already in use)
		panic("Failed to start server: " + err.Error())
	}
}