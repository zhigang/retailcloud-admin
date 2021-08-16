package config

// Config for Application.
type Config struct {
	Service struct {
		Address   string `default:":3000"`
		BasicAuth struct {
			Username string `default:"admin"`
			Password string `default:"admin"`
		}
	}

	Aksk struct {
		RegionID        string `required:"true" env:"AKSK_REGIONID"`
		AccessKeyID     string `required:"true" env:"AKSK_ACCESSKEYID"`
		AccessKeySecret string `required:"true" env:"AKSK_ACCESSKEYSECRET"`
	}

	Log struct {
		Level string `default:"debug"`
	}
}
