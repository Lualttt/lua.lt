package static

import "embed"

//go:embed css/*.css images/*.png javascript/*.js
var StaticContent embed.FS
