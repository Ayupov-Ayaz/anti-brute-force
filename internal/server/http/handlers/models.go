package handlers

type Auth struct {
	Login string `json:"login" validate:"required"`
	Pass  string `json:"password" validate:"required"`
	IP    string `json:"ip" validate:"required"`
}

type IP struct {
	IP   string `json:"ip" validate:"required"`
	Mask string `json:"mask" validate:"required"`
}

type Response struct {
	Ok bool `json:"ok"`
}
