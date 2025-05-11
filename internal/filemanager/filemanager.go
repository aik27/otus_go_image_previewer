package filemanager

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
)

func SaveFile(path string, body []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			slog.Error(fmt.Sprintf("unable to close os.File: %s", closeErr))
		}
	}(file)

	reader := bytes.NewReader(body)

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error read file: %w", err)
	}

	return content, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}

func GetFileNameByURL(url string) string {
	ext := path.Ext(url)
	hash := md5.Sum([]byte(url)) //nolint:gosec
	return fmt.Sprintf("%s%s", hex.EncodeToString(hash[:]), ext)
}
