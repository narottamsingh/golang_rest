package responses

import "github.com/gofiber/fiber"

type EmployeeResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}
