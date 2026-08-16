package workspace

import (
	"path"
)

func isAllowedTextFile(toolPath string) bool {
	basePath := path.Base(toolPath)
	switch basePath {
	case "go.mod", "go.sum":
		return true
	}

	suffix := path.Ext(basePath)
	switch suffix {
	case ".go", ".md", ".txt", ".json", ".yaml", ".yml", "toml":
		return true
	}

	return false
}
