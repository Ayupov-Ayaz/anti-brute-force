package handlers

type Auth struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
	IP    string `json:"ip"`
}

type IP struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}
