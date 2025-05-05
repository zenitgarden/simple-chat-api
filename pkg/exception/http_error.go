package exception

import "github.com/gofiber/fiber/v2"

type HTTPError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func JSON(c *fiber.Ctx, err error) error {
	if httpErr, ok := err.(*HTTPError); ok {
		return c.Status(httpErr.StatusCode).JSON(httpErr)
	}
	// fallback for generic errors
	return c.Status(fiber.StatusInternalServerError).JSON(HTTPError{
		StatusCode: fiber.StatusInternalServerError,
		Message:    "Internal Server Error",
		Data:       nil,
	})
}
