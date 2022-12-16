package httpsvc

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/irvankadhafi/articles-go/article-service/utils"
	"github.com/labstack/echo"
	"net/http"
)

var (
	ErrInvalidArgument = echo.NewHTTPError(http.StatusBadRequest, "invalid argument")
	ErrInternal        = echo.NewHTTPError(http.StatusInternalServerError, "internal system error")
)

// httpValidationOrInternalErr return valdiation or internal error
func httpValidationOrInternalErr(err error) error {
	switch t := err.(type) {
	case validator.ValidationErrors:
		_ = t
		errVal := err.(validator.ValidationErrors)

		fields := map[string]interface{}{}
		for _, ve := range errVal {
			fields[ve.Field()] = fmt.Sprintf("Failed on the '%s' tag", ve.Tag())
		}

		return echo.NewHTTPError(http.StatusBadRequest, utils.Dump(fields))
	default:
		return ErrInternal
	}
}
