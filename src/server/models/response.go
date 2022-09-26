package models

type ValidationErrorResponse struct {
	Message string `json:"message"`
	Errors  Errors `json:"errors"`
}

type Errors struct {
	AdditionalProperties string `json:"additional_properties"`
}

type PersonRequest struct {
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type PersonResponse struct {
	Id      int32  `json:"id"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
