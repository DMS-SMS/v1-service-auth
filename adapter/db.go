package adapter

type DBConfig struct {
	Dialect string `json:"dialect" validate:"required"`
	Host    string `json:"host" validate:"required"`
	Port 	int	   `json:"port" validate:"required"`
	User    string `json:"user" validate:"required"`
	DB		string `json:"db" validate:"required"`
}

