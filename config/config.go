package config

type ToolConfig struct {
	Proxy       string `yaml:"proxy" mapstructure:"proxy"`
	ProxyEnable bool   `yaml:"proxy_enable" mapstructure:"proxy_enable"`
	IsTableShow bool   `yaml:"is_table_show" mapstructure:"is_table_show"`
	IsApiMode   bool   `yaml:"is_api_mode" mapstructure:"is_api_mode"`
	Biu         struct {
		Api      string   `yaml:"api" mapstructure:"api"`
		Key      string   `yaml:"key" mapstructure:"key"`
		Port     string   `yaml:"port" mapstructure:"port"`
		IsPublic bool     `yaml:"is-public" mapstructure:"is-public"`
		Tags     []string `yaml:"tags" mapstructure:"tags"`
	}
	Api struct {
		Server  string `yaml:"server" mapstructure:"server"`
		Mongodb string `yaml:"mongodb" mapstructure:"mongodb"`
		Redis   string `yaml:"redis" mapstructure:"redis"`
	}
	Cookies struct {
		Aldzs      string `yaml:"aldzs" mapstructure:"aldzs"`
		Xlb        string `yaml:"xlb" mapstructure:"xlb"`
		Aiqicha    string `yaml:"aiqicha" mapstructure:"aiqicha"`
		Tianyancha string `yaml:"tianyancha" mapstructure:"tianyancha"`
		Tycid      string `yaml:"tycid" mapstructure:"tycid"`
		Qichacha   string `yaml:"qichacha" mapstructure:"qichacha"`
		QiMai      string `yaml:"qimai" mapstructure:"qimai"`
		ChinaZ     string `yaml:"chinaz" mapstructure:"chinaz"`
		Veryvp     string `yaml:"veryvp" mapstructure:"veryvp"`
	}
}
