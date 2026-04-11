package response

type Response struct {
	Message          string                 `json:"message"`
	Data             any                    `json:"data,omitempty"`
	Err              string                 `json:"error,omitempty"`
	ValidationErrors []ValidationFieldError `json:"validation_errors,omitempty"`
}

type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type paginationResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
}

type Result struct {
	Success          bool
	HttpStatus       int
	Message          string
	Data             any
	Error            error
	ValidationErrors []ValidationFieldError
}

func (r Result) IsError() bool {
	return !r.Success
}
