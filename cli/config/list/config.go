package listcfg

type IPList struct {
	BlackListAddr string `envconfig:"BLACK_LIST_ADDR" validate:"required"`
	WhiteListAddr string `envconfig:"WHITE_LIST_ADDR" validate:"required"`
}
