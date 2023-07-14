package static

import "embed"

//go:embed index.html
var HTML embed.FS

//go:embed webauthn.js webauthn.min.js
var JS embed.FS
