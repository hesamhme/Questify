package presenter

import "github.com/gofiber/fiber/v2"

// Response represents a common HTTP response structure.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse creates a standard success response.
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse creates a standard error response.
func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
	}
}

// WriteResponse sends the response to the client.
func WriteResponse(c *fiber.Ctx, status int, resp Response) error {
	return c.Status(status).JSON(resp)
}
