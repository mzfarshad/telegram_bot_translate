package contactus

import (
	"fmt"
	"os"
)

func ReadContactUsFile(path string) (string, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", path, err)
	}

	return string(content), nil
}
