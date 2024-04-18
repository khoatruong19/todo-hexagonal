package utils

import (
	"fmt"
	"net/http"
	"todo-hexagonal/internal/types"

	"todo-hexagonal/internal/adapters/primary/web/httperror"
	authLayout "todo-hexagonal/internal/adapters/primary/web/views/layout/auth"
	mainLayout "todo-hexagonal/internal/adapters/primary/web/views/layout/main"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
)

func RenderPage(w http.ResponseWriter, r *http.Request, layout string, page templ.Component, pageTitle string) {

	layouts := map[string]func(contents templ.Component, title string) templ.Component{
		"auth": authLayout.AuthLayout,
		"main": mainLayout.Layout,
	}

	err := layouts[layout](page, pageTitle).Render(r.Context(), w)
	if err != nil {
		httperror.RenderErrorResponse(w)
		return
	}
}

func RenderComponent(w http.ResponseWriter, r *http.Request, component templ.Component) {
	err := component.Render(r.Context(), w)
	if err != nil {
		httperror.RenderErrorResponse(w)
		return
	}
}

func errorMessageForTag(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Invalid email"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be no longer than %s characters", field, param)
	}

	return ""
}

func ValidationInput(data interface{}) *[]types.ValidationError {
	validate := validator.New()
	validationErrors := []types.ValidationError{}

	err := validate.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}

		for _, err := range err.(validator.ValidationErrors) {
			e := types.ValidationError{
				Field:   err.Field(),
				Message: errorMessageForTag(err.Field(), err.Tag(), err.Param()),
			}

			validationErrors = append(validationErrors, e)
		}

	}

	if len(validationErrors) > 0 {
		return &validationErrors
	}

	return nil
}

func NewValidationError(field, message string) types.ValidationError {
	return types.ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewValidationErrors(errors ...types.ValidationError) []types.ValidationError {
	return errors
}
