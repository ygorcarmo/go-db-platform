package models

type Response struct {
	Message   string
	IsSuccess bool
}

func CreateResponse(m string, i bool) Response {
	return Response{Message: m, IsSuccess: i}
}
