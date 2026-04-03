package configs

type ChatServer struct {
	BaseUrl   string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	AuthToken string `env:"AUTH_TOKEN" envDefault:"auth-token"`
	Model     string `env:"MODEL" envDefault:"openai"`
}
