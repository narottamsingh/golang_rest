package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func EmployeeRoute(app *fiber.App) {
	app.Post("/employee", controllers.CreateEmployee)          // create new employee
	app.Get("/employee", controllers.GetEmployeeList)          // get all employee
	app.Get("/employee/:empId", controllers.GetEmployee)       // get emplpyee by id
	app.Put("/employee/:empId", controllers.EditEmployee)      // edit empliyee
	app.Delete("/employee/:empId", controllers.DeleteEmployee) // delete employee

}
