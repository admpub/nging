package middleware

var Middlewares []interface{}

func Use(m ...interface{}) {
	Middlewares = append(Middlewares, m...)
}
