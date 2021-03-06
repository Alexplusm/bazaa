package httputils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// ParseContentType parse content type
// TODO: may be stdlib func this same functional exsists?
// TODO: see "mime.ParseMediaType"
func ParseContentType(ctx echo.Context) string {
	// fmt.Printf("HEADER %+v\n", ctx.Request().Header)
	header := ctx.Request().Header
	contentTypeRaw, ok := header["Content-Type"]
	if !ok {
		return "content type undefined"
	}
	contentType := strings.Split(contentTypeRaw[0], ";")[0]

	return contentType
}

// "application/json"

// ContentTypeMiddleware check content-type
func ContentTypeMiddleware(ctx echo.Context, expectedContentType string) error {
	fmt.Println("Middleware") // logger
	ctxContentType := ParseContentType(ctx)

	if ctxContentType != expectedContentType {
		// TODO: test
		ctx.JSON(http.StatusOK, BuildBadRequestErrorResponse())
		return fmt.Errorf("Error content-type")
	}
	return nil
}
