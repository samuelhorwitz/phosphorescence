package session

var (
	cookieDomain string
	isProduction bool
)

type Config struct {
	CookieDomain string
	IsProduction bool
}

func Initialize(cfg *Config) {
	cookieDomain = cfg.CookieDomain
	isProduction = cfg.IsProduction
	initializeReaper()
}
