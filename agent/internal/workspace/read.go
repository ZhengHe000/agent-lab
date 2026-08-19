package workspace

import (
	"fmt"
	"io"
)

const maxReadBytes = 40_000

func (w *Workspace) ReadTextFile(input string) (string, error) {
	toolPath, err := validateToolPath(input)
	if err != nil {
		return "", fmt.Errorf("validate path %v: %w", input, err)
	}

	if !isAllowedTextFile(toolPath) {
		return "", fmt.Errorf("validate path %v: %w", toolPath, ErrUnsupportedFileType)
	}

	localPath, err := localizeToolPath(toolPath)
	if err != nil {
		return "", err
	}

	if err = w.rejectSymlinkPath(localPath); err != nil {
		return "", err
	}

	file, err := w.root.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("open file %v: %w", localPath, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file %v: %w", toolPath, err)
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("validate file %v: not a regular file", toolPath)
	}

	reader := io.LimitReader(file, int64(maxReadBytes+1))

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read file %v: %w", toolPath, err)
	}

	if len(data) > maxReadBytes {
		return "", fmt.Errorf("read file %v: %w", toolPath, ErrFileTooLarge)
	}

	return string(data), nil
}
