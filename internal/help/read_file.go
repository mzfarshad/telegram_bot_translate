package help

import (
	"fmt"
	"os"
)

func ReadHelpFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading help file %s: %v", path, err)
	}
	return string(content), nil
}
