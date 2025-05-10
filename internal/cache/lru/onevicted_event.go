package lrucache

import (
	"context"
	"log/slog"

	"github.com/aik27/otus_go_image_previewer/internal/filemanager"
)

func OnEvictedEvent(item *CacheItem) {
	fm := filemanager.NewFileManager(context.Background())
	if imageItem, ok := item.Value.(ImageItem); ok {
		err := fm.DeleteFile(imageItem.FilePath)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
