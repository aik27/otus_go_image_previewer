package events

import (
	"context"
	"log/slog"

	lrucache "github.com/aik27/otus_go_image_previewer/internal/cache/lru"
	"github.com/aik27/otus_go_image_previewer/internal/filemanager"
)

func OnEvicted(item *lrucache.CacheItem) {
	fm := filemanager.NewFileManager(context.Background())
	if imageItem, ok := item.Value.(lrucache.ImageItem); ok {
		err := fm.DeleteFile(imageItem.FilePath)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
