package model

type AccessLogLite struct {
	Date   string
	OS     string
	Brower string
	Region string
	Type   string

	Version string
	User    string
	Method  string
	Scheme  string
	Host    string
	URI     string

	Referer string

	BodyBytes  uint64
	Elapsed    float64
	StatusCode uint
	UserAgent  string
}
