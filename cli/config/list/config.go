package listcfg

type IPList struct {
	BlackListAddr string `mapstructure:"blacklist_addr" validate:"required"`
	WhiteListAddr string `mapstructure:"whitelist_addr" validate:"required"`
}
