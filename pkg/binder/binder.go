package binder

import (
    "context"
    "fmt"
    "net/http"
    "reflect"

    "github.com/creasty/defaults"
    "github.com/labstack/echo"
    "gopkg.in/go-playground/mold.v2"
    "gopkg.in/go-playground/mold.v2/modifiers"
    "gopkg.in/go-playground/validator.v9"
)

// Binder is a custom struct that implements the Echo Binder interface. It binds
// to a struct, uses mold to clean up the params, and validator to validate
// them.
type Binder struct {
    binder   *echo.DefaultBinder
    conform  *mold.Transformer
    validate *validator.Validate
}

// New initializes a new Binder instance with the appropriate validation
// functions registered.
func New() *Binder {
    binder := &echo.DefaultBinder{} // echo's default binder
    conform := modifiers.New()
    validate := validator.New()

    return &Binder{binder, conform, validate}
}

// Bind binds, modifies, and validates payloads against the given struct.
func (b *Binder) Bind(i interface{}, c echo.Context) error {
    // Extract values from the request in the echo.Context into our interace i.
    // Note that we are still using the default echo binder to extract the data
    // from the echo.Context but we are doing some extra validation in the
    // process.
    if err := b.binder.Bind(i, c); err != nil {
        return err
    }

    // Modify values based on the struct tags in our interface. Most likely this
    // is trimming whitespace from values.
    if err := b.conform.Struct(context.Background(), i); err != nil {
        return err
    }

    // Set any default values that we have set.
    if err := defaults.Set(i); err != nil {
        return err
    }

    // Validate that the values on our struct are valid. If there are any errors,
    // format the first error and return it as an HTTP error.
    if err := b.validate.Struct(i); err != nil {
        errs := err.(validator.ValidationErrors)
        msg := format(errs[0])
        return echo.NewHTTPError(http.StatusUnprocessableEntity, msg)
    }

    return nil
}

func format(err validator.FieldError) string {
    if err.Kind() == reflect.Int {
        switch err.Tag() {
        case "max":
            return fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
        case "min":
            return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
        }
    } else if err.Kind() == reflect.String {
        switch err.Tag() {
        case "max":
            return fmt.Sprintf("%s length must be less than or equal to %s characters long", err.Field(), err.Param())
        case "min":
            return fmt.Sprintf("%s length must be at least %s characters long", err.Field(), err.Param())
        }
    }

    if err.Tag() == "required" {
        return fmt.Sprintf("%s is required", err.Field())
    }

    return err.Param()
}
