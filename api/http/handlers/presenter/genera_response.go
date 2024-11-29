package presenter

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Success: true,
	}
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) SetError(err error) *Response {
	r.Success = false
	r.Error = err.Error()
	return r
}

func Send(c *fiber.Ctx, status int, response *Response) error {
	return c.Status(status).JSON(response)
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
	return Send(c, fiber.StatusOK, NewResponse().SetMessage(message).SetData(data))
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return Send(c, fiber.StatusCreated, NewResponse().SetMessage(message).SetData(data))
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func BadRequest(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusBadRequest, NewResponse().SetError(err))
}

func Unauthorized(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusUnauthorized, NewResponse().SetError(err))
}

func Forbidden(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusForbidden, NewResponse().SetError(err))
}

func Conflict(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusConflict, NewResponse().SetError(err))
}

func NotFound(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusNotFound, NewResponse().SetError(err))
}

func InternalServerError(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusInternalServerError, NewResponse().SetError(err))
}