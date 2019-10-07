package codec

type Codec interface {
	Encode(raw string, keys ...string) string
	Decode(encrypted string, keys ...string) string
}
