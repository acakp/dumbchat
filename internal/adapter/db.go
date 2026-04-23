package adapter

type Config struct {
	DBDriver   string `env:"DB" envDefault:"sqlite"`
	PGHost     string `env:"PGHOST"`
	PGPort     string `env:"PGPORT"`
	PGDBName   string `env:"PGDBNAME"`
	PGUser     string `env:"PGUSER"`
	PGPassword string `env:"PGPASSWORD"`
}
