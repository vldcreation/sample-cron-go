package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/vldcreation/sample-cron-go/internal/config"
)

var (
	ErrNotFound     = fmt.Errorf("storage object not found")
	conf            = config.NewAppConfig()
	Test5Seconds    = time.Second * 5    // 5 seconds, for testing
	Test10Seconds   = time.Second * 10   // 10 seconds, for testing
	Test20Seconds   = time.Second * 20   // 20 seconds, for testing
	TestDuration    = time.Minute * 1    // 1 minute, for testing
	DefaultDuration = time.Hour * 24     // 1 days
	MaxDuration     = time.Hour * 24 * 7 // 7 days, max expiry for presigned URLs
)

type PresignUrlInfoS3 struct {
	X_AMZ_ALGORITHM     string `json:"X-Amz-Algorithm"`
	X_AMZ_CREDENTIAL    string `json:"X-Amz-Credential"`
	X_AMZ_DATE          string `json:"X-Amz-Date"`
	X_AMZ_EXPIRES       string `json:"X-Amz-Expires"`
	X_AMZ_SIGNEDHEADERS string `json:"X-Amz-SignedHeaders"`
	X_AMZ_SIGNATURE     string `json:"X-Amz-Signature"`
}

type PresigneUrlInfoGCS struct {
	GoogleAccessId string `json:"GoogleAccessId"`
	Expires        string `json:"Expires"`
	Signature      string `json:"Signature"`
}

type Storage interface {
	// Put creates or overwrites an object in the storage system.
	// If contentType is blank, the default for the chosen storage implementation is used.
	Put(ctx context.Context, parent, name string, contents []byte, cacheAble bool, contentType string) error

	// FPutObject creates or overwrites an object in the storage system from filepath.
	FPut(ctx context.Context, parent, name, filePath string, cacheAble bool, contentType string) error

	// Delete deletes an object or does nothing if the object doesn't exist.
	Delete(ctx context.Context, parent, bame string) error

	// Get fetches the object's contents.
	Get(ctx context.Context, parent, name string) ([]byte, error)

	// PresignURL returns a presigned URL for the object with replace versioning file.
	ReSignedURLWithReplace(ctx context.Context, parent, object string) (string, error)

	// PresignURL returns a presigned URL for the object.
	ReSignedURL(ctx context.Context, parent, object, existingUrl string) (string, error)
}
