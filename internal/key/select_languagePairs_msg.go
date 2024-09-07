package key

import (
	"fmt"
	"os"
)

func ReadSelectLanguagePairsMsg(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading select language pairs message file %s: %v", path, err)
	}

	return string(content), nil
}
