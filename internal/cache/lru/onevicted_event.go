package lrucache

import (
	"log/slog"

	"github.com/aik27/otus_go_image_previewer/internal/filemanager"
)

func OnEvictedEvent(item *CacheItem) {
	if imageItem, ok := item.Value.(ImageItem); ok {
		err := filemanager.DeleteFile(imageItem.FilePath)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
