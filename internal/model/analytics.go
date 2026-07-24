package model

type CalculateRequest struct {
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
	From        string `form:"from"`
	To          string `form:"to"`
}

type CalculateResponse struct {
	Total int `json:"total"`
}