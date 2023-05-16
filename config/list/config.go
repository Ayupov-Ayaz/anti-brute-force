package listcfg

type IPList struct {
	BlackListAddr string `mapstructure:"blacklist_addr"`
	WhiteListAddr string `mapstructure:"whitelist_addr"`
}
