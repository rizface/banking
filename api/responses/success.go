package responses

import "github.com/gofiber/fiber/v2"

type TheResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TheResponseUpload struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type TheResponseWithMeta struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

func ReturnTheResponse(c *fiber.Ctx, se bool, sc int, m string, dt interface{}) error {
	tr := TheResponse{m, dt}

	return c.Status(sc).JSON(tr)
}

func ReturnTheResponseMeta(c *fiber.Ctx, se bool, sc int, m string, dt interface{}, meta Meta) error {
	tr := TheResponseWithMeta{m, dt, meta}

	return c.Status(sc).JSON(tr)
}

func ReturnTheResponseUpload(c *fiber.Ctx, se bool, sc int, m string, f string) error {
	tr := TheResponseUpload{m, f}

	return c.Status(sc).JSON(tr)
}
