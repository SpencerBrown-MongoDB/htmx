package web

import "embed"

// Static content as string variables

//go:embed static template
var asset embed.FS
