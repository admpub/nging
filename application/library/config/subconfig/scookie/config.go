package scookie

type Config struct {
	Domain   string `json:"domain"`
	MaxAge   int    `json:"maxAge"`
	Path     string `json:"path"`
	HttpOnly bool   `json:"httpOnly"`
	HashKey  string `json:"hashKey"`
	BlockKey string `json:"blockKey"`
}
