package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asimn7/two-cc-project/models"
	"github.com/asimn7/two-cc-project/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

//MARK: Repository
type Repository struct{
	DB *gorm.DB
}

//MARK: Routes
func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_foods", r.CreateFood)
	api.Delete("/delete_foods/:id", r.DeleteFood)
	api.Get("/get_foods/:id", r.GetFoodById)
	api.Get("/foods", r.GetFoods)
}

//MARK: Food struct
type Food struct{
	Seller			string		`json:"seller"`
	Buyer			string		`json:"buyer"`
	Publisher		string		`json:"publisher"`
}

//MARK: CreateBook
func (r *Repository) CreateFood(c *fiber.Ctx) error{
	food := Food{}

	err := c.BodyParser(&food)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
		return err
	}

	err = r.DB.Create(&food).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"an error occurred while creating the food"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"food created successfully", "data":food})

	return nil

}

//MARK: DeleteFood
func (r *Repository) DeleteFood (c *fiber.Ctx) error{
	foodModel := models.Foods{}

	id := c.Params("id")

	if id == ""{
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(foodModel, id)

	if err.Error != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not delete food"})
		return err.Error
	}
	
	c.Status(http.StatusOK).JSON(&fiber.Map{"message":"food deleted successfully"})

	return nil
}

//MARK: GetBookById
func (r *Repository) GetFoodById(c *fiber.Ctx) error{
	foodModel := models.Foods{}
	id := c.Params("id")

	if id == ""{
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"id cannot be empty"})
		return nil
	}

	err := r.DB.First(&foodModel, id)

	if err.Error != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not find food"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "food found successfully", "data":foodModel})

	return nil
}

//MARK: Get All Books
func (r *Repository) GetFoods(c *fiber.Ctx) error{
	foodModels := &[]models.Foods{}

	err := r.DB.Find(foodModels).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"an error occurred while fetching the foods"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"foods fetched successfully", "data":foodModels})
	return nil
}

//MARK: Main
func main(){
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	 //MARK: Database connection
	 config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User: os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSL_MODE"),
		DBName: os.Getenv("DB_NAME"),
	 }
	 
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = models.MigrateFoods(db)

	if err != nil {
		log.Fatal("Error migrating database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8083")
}