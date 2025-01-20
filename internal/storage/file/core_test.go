package file

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileManager(t *testing.T) {
	t.Run("success init file-manager", func(t *testing.T) {
		dir, err := os.MkdirTemp("./", "tempfile")
		if err != nil {
			t.Fatalf("Не удалось создать временную директорию: %v", err)
		}
		defer os.RemoveAll(dir)

		fm, err := NewManager(dir)
		if err != nil {
			t.Fatalf("NewManager вернул ошибку: %v", err)
		}
		defer func() {
			if fm.file != nil {
				fm.file.Close()
			}
		}()

		if fm.file == nil {
			t.Error("ожидалось, что file не равен nil")
		}
		if fm.producer == nil {
			t.Error("ожидалось, что producer не равен nil")
		}
		if fm.consumer == nil {
			t.Error("ожидалось, что consumer не равен nil")
		}
	})
	t.Run("error init file-manager", func(t *testing.T) {
		random := strconv.Itoa(rand.Int())
		_, err := NewManager("fakedir" + random)
		assert.Error(t, err, errors.New("no such file or directory"))
	})
}
