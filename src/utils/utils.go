package utils

import (
	"encoding/base64"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func ParseBase64(uuid string) (string, error) {
	rawDecodedText, err := base64.StdEncoding.DecodeString(uuid)
	return string(rawDecodedText), err
}

func RenderResponse(c *fiber.Ctx, component templ.Component) error {
	c.Response().Header.SetContentType("text/html")
	return component.Render(c.UserContext(), c.Context().Response.BodyWriter())
}