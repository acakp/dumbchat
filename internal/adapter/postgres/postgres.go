package postgres

type Config struct {
	PGHost     string `env:"PGHOST"`
	PGPort     string `env:"PGPORT"`
	PGDBName   string `env:"PGDBNAME"`
	PGUser     string `env:"PGUSER"`
	PGPassword string `env:"PGPASSWORD"`
}
