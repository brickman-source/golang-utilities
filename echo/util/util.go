package util

import (
	"errors"

	"github.com/labstack/echo/v4"

	httpUtil "github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
)

func ParseJSON(ctx echo.Context, v interface{}) error {
	data := httpUtil.GetRequestBody(ctx.Request())
	if data == nil {
		return errors.New("bad request")
	}
	return json.Unmarshal(data, v)
}

func WriteJSON(ctx echo.Context, code int, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ctx.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, data)
}

func WriteData(ctx echo.Context, code int, contentType string, data []byte) error {
	return ctx.Blob(code, contentType, data)
}
