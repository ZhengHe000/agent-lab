package workspace

import (
	"io/fs"
	"strings"
)

func validateToolPath(input string) (string, error) {
	if input == "." {
		return "", ErrInvalidPath
	}

	if !fs.ValidPath(input) {
		return "", ErrInvalidPath
	}

	if strings.ContainsAny(input, "\\:\x00") {
		return "", ErrInvalidPath
	}

	return input, nil
}
