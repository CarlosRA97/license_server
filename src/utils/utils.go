package utils

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func ParseBase64(uuid string) string {
	rawDecodedText, err := base64.StdEncoding.DecodeString(uuid)
	if err != nil {
		panic(err)
	}
	return string(rawDecodedText)
}

func Render(component templ.Component) string {
	html := new(bytes.Buffer)
	component.Render(context.Background(), html)
	return html.String()
}

func RenderResponse(c *fiber.Ctx, component templ.Component) error {
	c.Response().Header.SetContentType("text/html")
	return component.Render(c.UserContext(), c.Context().Response.BodyWriter())
}