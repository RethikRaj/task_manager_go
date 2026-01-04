package dto

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
