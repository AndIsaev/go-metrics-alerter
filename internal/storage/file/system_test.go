package file

import (
	"os"
	"testing"
)

func TestCreateDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file, _ := os.CreateTemp(tempDir, "tempfile")
	defer file.Close()
	fm := &FileManager{file: file}

	t.Run("succes", func(t *testing.T) {
		newDirPath := tempDir + "/newdir"
		if err := fm.CreateDir(newDirPath); err != nil {
			t.Errorf("Ожидалось, что директория создана успешно, но возникла ошибка: %v", err)
		}
		if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
			t.Errorf("Директория %s должна существовать, но её нет", newDirPath)
		}
	})
	t.Run("dir already exist", func(t *testing.T) {
		newDirPath := tempDir + "/newdir"
		if err := fm.CreateDir(newDirPath); err != nil {
			t.Errorf("Ожидалось, что метод завершится успешно, несмотря на существование директории, но возникла ошибка: %v", err)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		if err := fm.CreateDir("/root/forbidden"); err == nil {
			t.Errorf("Ожидалась ошибка при создании директории в защищённой области, но её не было")
		}
	})
}
