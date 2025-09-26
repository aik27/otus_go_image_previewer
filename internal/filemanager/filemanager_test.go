package filemanager

import (
	"bytes"
	"os"
	"path"
	"testing"
)

func TestFilemanager(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		filePath string
		body     []byte
	}{
		{
			name:     "test1",
			filePath: path.Join(os.TempDir(), "test1.txt"),
			body:     []byte("test1"),
		},
		{
			name:     "test2",
			filePath: path.Join(os.TempDir(), "test2.txt"),
			body:     []byte("test2"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := SaveFile(tc.filePath, tc.body)
			if err != nil {
				t.Fatalf("failed to save file: %v", err)
			}

			content, err := ReadFile(tc.filePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if !bytes.Equal(content, tc.body) {
				t.Fatalf("file content mismatch: got %s, want %s", content, tc.body)
			}

			err = DeleteFile(tc.filePath)
			if err != nil {
				t.Fatalf("failed to delete file: %v", err)
			}
		})
	}
}

func TestGetFileNameByURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "test1",
			filePath: "http://example.com/test1.jpg",
			expected: "34f8ed9f478dc7921389fabc6d5ae198.jpg",
		},
		{
			name:     "test2",
			filePath: "https://example.com/test2.png",
			expected: "950dc12dd0ce7ecfb604051371786269.png",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := GetFileNameByURL(tc.filePath)
			if result != tc.expected {
				t.Fatalf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}
