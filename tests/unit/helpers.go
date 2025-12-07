package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/api-service/config"
	"github.com/api-service/logger"
	"github.com/api-service/storage"
)

const (
	maxRetryCount = 3
	retryTimeout  = 5 * time.Second
)

func GetConfig(ctx context.Context) (*config.Config, error) {
	return config.GetConfig(ctx)
}

func GetStorage(ctx context.Context) (storage.Storage, error) {
	return storage.GetStorage(ctx)
}

func GetLogger(ctx context.Context) logger.Logger {
	return logger.GetLogger(ctx)
}

func RetryWithTimeout(ctx context.Context, f func() error, timeout time.Duration) error {
	for attempt := 0; attempt <= maxRetryCount; attempt++ {
		if err := f(); err == nil {
			return nil
		}
		if attempt >= maxRetryCount {
			return err
		}
		select {
		case <-ctx.Done():
			return errors.New("context cancelled")
		case <-time.After(timeout):
		}
	}
	return errors.New("max retry count exceeded")
}

func HandleHTTPError(w http.ResponseWriter, status int, err error) {
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), status)
		return
	}
}

func FormatError(err error) string {
	return fmt.Sprintf("Error: %v", err)
}