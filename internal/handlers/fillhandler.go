package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"time"

	lrucache "github.com/aik27/otus_go_image_previewer/internal/cache/lru"
	"github.com/aik27/otus_go_image_previewer/internal/config"
	"github.com/aik27/otus_go_image_previewer/internal/filemanager"
	"github.com/aik27/otus_go_image_previewer/internal/imageprocessor"
	"github.com/aik27/otus_go_image_previewer/internal/proxy"
	"github.com/go-chi/chi/v5"
)

func FillHandler(
	ctx context.Context,
	cfg *config.Config,
	cache lrucache.Cache,
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	proxyClient := proxy.NewClient(ctx)

	width := chi.URLParam(r, "width")
	height := chi.URLParam(r, "height")
	imgURL := chi.URLParam(r, "*")

	wInt, err := strconv.Atoi(width)
	if err != nil {
		http.Error(w, "Invalid width value", http.StatusBadRequest)
		return
	}

	hInt, err := strconv.Atoi(height)
	if err != nil {
		http.Error(w, "Invalid height value", http.StatusBadRequest)
		return
	}

	ext := path.Ext(imgURL)
	if ext != ".jpg" {
		http.Error(w, "Invalid file extension. JPG supported only", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("%s_%s_%s", width, height, imgURL)
	cached, ok := cache.Get(lrucache.Key(cacheKey))

	if ok {
		slog.Debug(fmt.Sprintf("Hit to cache: %s", cacheKey))

		filePath := cached.(lrucache.ImageItem).FilePath
		file, err := filemanager.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Failed to read file from cache", http.StatusInternalServerError)
			slog.Error(fmt.Sprintf("Failed to read file from cache: %s: %s", filePath, err.Error()))
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(file)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			slog.Error(fmt.Sprintf("Failed to write response: %s", err.Error()))
			return
		}

		return
	}

	image, err := proxyClient.FetchFile(imgURL, r)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Failed to fetch image: %s", err.Error()))
		return
	}

	image, err = imageprocessor.Resize(image, wInt, hInt)
	if err != nil {
		http.Error(w, "Failed to modify image.", http.StatusInternalServerError)
		return
	}

	savePath := fmt.Sprintf("%s/%s_%s_%s", cfg.Cache.Dir, width, height, filemanager.GetFileNameByURL(imgURL))

	err = filemanager.SaveFile(savePath, image)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Failed to save image: %s", err.Error()))
		return
	}

	cached = lrucache.ImageItem{
		FilePath:    savePath,
		Width:       wInt,
		Height:      hInt,
		OriginalURL: imgURL,
	}

	_ = cache.Set(lrucache.Key(cacheKey), cached)

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(image)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to write response: %s", err.Error()))
		return
	}
}
