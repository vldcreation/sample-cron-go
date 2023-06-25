package utils_test

import (
	"testing"

	"github.com/vldcreation/sample-cron-go/internal/utils"
)

func TestParseFile(t *testing.T) {
	t.Run("should return error when filename is invalid", func(t *testing.T) {
		_, _, err := utils.ParseFile("file")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("should return name and extension when filename is valid", func(t *testing.T) {
		name, ext, err := utils.ParseFile("./test_data/file.txt")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if name != "file" {
			t.Errorf("expected file, got %v", name)
		}

		if ext != "txt" {
			t.Errorf("expected txt, got %v", ext)
		}
	})
}
