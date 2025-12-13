package assets

import "embed"

//go:embed statics/* assets/*
var Assets embed.FS
