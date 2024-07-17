package handlers

type BaseRequest struct {
	IP    string `json:"ip" validate:"required"`
	Login string `json:"login" validate:"required"`
}

type CheckAuthRequest struct {
	BaseRequest
	Pass string `json:"password" validate:"required"`
}

type ResetRequest struct {
}

type IP struct {
	IPNet string `json:"ip_net" validate:"required"`
}

type Response struct {
	Ok bool `json:"ok"`
}
