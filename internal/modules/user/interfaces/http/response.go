package http

type MeResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	Status   string `json:"status"`
}
