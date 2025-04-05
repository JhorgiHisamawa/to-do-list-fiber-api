package main

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

	_ "github.com/lib/pq"
)

// Activity model
type Activity struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title" validate:"required"`
	Category     string    `json:"category" validate:"required,oneof=TASK EVENT"`
	Description  string    `json:"description" validate:"required"`
	ActivityDate time.Time `json:"activityDate" validate:"required"`
	Status       string    `json:"status" validate:"required,oneof=NEW 'ON PROGRESS' EXPIRED"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Set up a connection to the database
func initDB() (*sql.DB, error) {
	dns := "postgresql://<username>:<password>@<host>:<port>/<database>"
	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// Run the server
func main() {

	db, err := initDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	app := fiber.New()
	validate := validator.New()

	/*
		Define routes and logic handlers
	*/

	// Get all activities
	app.Get("/activities", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT * FROM activities")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		defer rows.Close()

		var activities []Activity
		for rows.Next() {
			var activity Activity
			if err := rows.Scan(&activity.ID, &activity.CreatedAt, &activity.Title, &activity.Category, &activity.Description, &activity.ActivityDate, &activity.Status); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			}
			activities = append(activities, activity)
		}

		return c.Status(fiber.StatusOK).JSON(activities)
	})

	// Create a new activity
	app.Post("/activities", func(c *fiber.Ctx) error {
		var activity Activity
		if err := c.BodyParser(&activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		if err = validate.Struct(activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		sql := "INSERT INTO activities (title, category, description, activity_date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		err := db.QueryRow(sql, activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status).Scan(&activity.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Success", "id": activity.ID})
	})

	// Update an existing activity
	app.Put("/activities/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var activity Activity
		if err := c.BodyParser(&activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		if err = validate.Struct(activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		sql := "UPDATE activities SET title = $1, category = $2, description = $3, activity_date = $4, status = $5 WHERE id = $6 RETURNING id"
		err := db.QueryRow(sql, activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status, id).Scan(&activity.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success"})
	})

	// Delete an activity
	app.Delete("/activities/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		sql := "DELETE FROM activities WHERE id = $1"
		_, err := db.Exec(sql, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success"})
	})

	// Start the server
	app.Listen(":3000")
}
