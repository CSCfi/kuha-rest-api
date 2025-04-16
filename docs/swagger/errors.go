package swagger

//General error response

// 400
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Error1 string `json:"parameter name" example:"full description of the error"`
}

// 401
type UnauthorizedResponse struct {
	Errors []UnauthorizedError `json:"errors"`
}

type UnauthorizedError struct {
	Error1 string `json:"error" example:"unauthorized"`
}

// 500
type InternalServerErrorResponse struct {
	Errors []InternalServerError `json:"errors"`
}

type InternalServerError struct {
	Error1 string `json:"error" example:"the server encountered a problem"`
}

// 422 - OURA, Polar, Suunto, Garmin
type OuraInvalidDateRange struct {
	Errors []UnprocessableEntityError `json:"errors"`
}

type UnprocessableEntityError struct {
	Error1 string `json:"error" example:"invalid date range"`
}
