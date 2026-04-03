package configs

import "time"

type Auth struct {
	JWTSecret     string        `env:"JWT_SECRET"           envDefault:"change-me-super-secret-jwt-key"`
	JWTAccessTTL  time.Duration `env:"JWT_ACCESS_TTL"       envDefault:"30m"`
	JWTRefreshTTL time.Duration `env:"JWT_REFRESH_TTL"      envDefault:"72h"`
}
