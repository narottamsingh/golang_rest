package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var employeeCollection *mongo.Collection = configs.GetCollection(configs.DB, "employee")
var validate = validator.New()

func CreateEmployee(e *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var employee models.Employee
	defer cancel() // Ensure the context is canceled when the function returns

	if err := e.BodyParser(&employee); err != nil {
		//return e.Status(http.StatusBadRequest).JSON(responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"data":    err.Error(),
		})
	}

	//emploee the validator library to validate required fields
	if err := validate.Struct(&employee); err != nil {
		//return e.Status(http.StatusBadRequest).JSON(responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"data":    err.Error(),
		})
	}
	newEmployee := models.Employee{
		Id:   primitive.NewObjectID(),
		Name: employee.Name,
		Age:  employee.Age,
	}

	result, err := employeeCollection.InsertOne(ctx, newEmployee)
	if err != nil {
		//return e.Status(http.StatusInternalServerError).JSON(responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"data":    err.Error(),
		})
	}

	//return e.Status(http.StatusCreated).JSON(responses.EmployeeResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
	return e.Status(http.StatusBadRequest).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"data":    result,
	})
}

func EditEmployee(e *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	employeeID := e.Params("empId")
	var employee models.Employee
	defer cancel() // Ensure the context is canceled when the function returns

	objId, _ := primitive.ObjectIDFromHex(employeeID)

	update := bson.M{"name": employee.Name, "age": employee.Age}

	result, err := employeeCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Some error occure while updating employee : ",
			"data":    err.Error(),
		})
	}

	//get updated user details
	var updatedEmployee models.Employee
	if result.MatchedCount == 1 {
		err := employeeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedEmployee)
		if err != nil {
			return e.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Some error occure while updating employee : ",
				"data":    err.Error(),
			})
		}
	}

	return e.Status(http.StatusBadRequest).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    updatedEmployee,
	})
}

func DeleteEmployee(e *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	employeeID := e.Params("empId")
	defer cancel() // Ensure the context is canceled when the function returns

	objId, _ := primitive.ObjectIDFromHex(employeeID)

	result, err := employeeCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Some error occure while deleting employee : " + employeeID,
			"data":    err.Error(),
		})
	}
	if result.DeletedCount < 1 {
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "The employee id does not exist : " + employeeID,
			"data":    err.Error(),
		})
	}

	return e.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Successfully deleted",
		"data":    "Done!",
	})
}

func GetEmployee(e *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var employee models.Employee
	employeeID := e.Params("empId")
	defer cancel() // Ensure the context is canceled when the function returns

	objId, _ := primitive.ObjectIDFromHex(employeeID)

	err := employeeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&employee)
	if err != nil {
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"data":    err.Error(),
		})
	}
	return e.Status(http.StatusBadRequest).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    employee,
	})

}

func GetEmployeeList(e *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var employees []models.Employee
	defer cancel() // Ensure the context is canceled when the function returns

	results, err := employeeCollection.Find(ctx, bson.M{})

	if err != nil {
		return e.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"data":    err.Error(),
		})
	}
	// reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleEmp models.Employee
		if err = results.Decode(&singleEmp); err != nil {
			return e.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"data":    err.Error(),
			})

		}
		employees = append(employees, singleEmp)
	}
	return e.Status(http.StatusBadRequest).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    employees,
	})

}
