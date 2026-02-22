package web

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var StaticFiles embed.FS

func (s *Server) handleIndex(c *fiber.Ctx) error {
	// Try to serve embedded static/index.html
	data, err := StaticFiles.ReadFile("static/index.html")
	if err != nil {
		// Fallback to file system
		return c.SendFile("./web/static/index.html")
	}
	return c.Type("html").Send(data)
}

func GetStaticFS() fs.FS {
	f, err := fs.Sub(StaticFiles, "static")
	if err != nil {
		return StaticFiles
	}
	return f
}
