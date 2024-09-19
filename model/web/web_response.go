package web

type WebResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success message"`
	Data    interface{} `json:"data"`
}
