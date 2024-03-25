package validate

import "github.com/go-playground/validator"
type ValidateError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
func ValidateRequest(requestData any) []ValidateError {
   
    validate := validator.New()
    err := validate.Struct(requestData)
    if err != nil {
        var error []ValidateError
        for _, err := range err.(validator.ValidationErrors) {
            errorMsg := ValidateError{
                Field: err.Field(),
                Message:  err.ActualTag(),
            } 
            error = append(error, errorMsg)
        }
        return error
    }
    return nil
}