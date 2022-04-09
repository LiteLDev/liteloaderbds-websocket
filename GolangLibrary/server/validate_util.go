package server

import (
	"BDSWebsocket/server/logger"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// validate is Types validator
// we share single instance of validate, it caches struct info
var validate *validator.Validate = validator.New()

func ValidateStructPrint(t interface{}) bool {
	var (
		err error
	)
	err = validate.Struct(t)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok && validationErrors != nil {

			logger.Error.Printf("%s Validate error: ", reflect.TypeOf(t).Name())
			for _, validateErr := range validationErrors {
				logger.Error.Printf(" -> %s", validateErr.Error())
			}
			return false
		}
		return true
	}
	return true
}
