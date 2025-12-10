package pkg

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	//Reason string `json:"detail,omitempty"`
	//Help    string `json:"help,omitempty"`
}

// id results
type IdResult[T any] struct {
	Id T `json:"id"`
}

/*
type idResult8 struct {
	Id int8 `json:"id"`
}
type idResult16 struct {
	Id int16 `json:"id"`
}
type idResult64 struct {
	Id int64 `json:"id"`
}
type idResultUUID struct {
	Id uuid.UUID `json:"id"`
}
*/
