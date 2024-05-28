package config

type Config struct {
	Token    Token
	Email    Email
	Redis    Redis
	Midtrans Midtrans
	MongoDB  MongoDB
	Server   Server
	Auth     Auth
	Google   Google
	Facebook Facebook
}

type Server struct {
	Host string
	Port string
}

type Token struct {
	Secret_Key string
}

type Email struct {
	Host     string
	Name     string
	Password string
}

type Redis struct {
	Addr string
	Pass string
}

type Midtrans struct {
	Key    string
	IsProd bool
}

type MongoDB struct {
	URI string
}

type Auth struct {
	Secret_Key          string
	MaxAge              string
	IsProd              string
	GoogleCallBackUrl   string
	FacebookCallBackUrl string
}

type Google struct {
	ClientID     string
	ClientSecret string
	ScopeEmail   string
	ScopeProfile string
	State        string
	TokenUrl     string
}
type Facebook struct {
	ClientID     string
	ClientSecret string
}
