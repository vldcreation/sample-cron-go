package storage

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/vldcreation/sample-cron-go/internal/utils"
	"google.golang.org/api/option"
)

// Compile-time check to verify implements interface.
var (
	_ Storage = (*GCS)(nil)
)

// gcs implements the Blob interface and provides the ability
// write files to Google Cloud Storage.
type GCS struct {
	client *storage.Client
}

// NewGCS creates a Google Cloud Storage Client
func NewGCS(ctx context.Context, cfgJsonFIle string) (Storage, error) {
	credOpt := option.WithCredentialsFile(cfgJsonFIle)

	client, err := storage.NewClient(ctx, credOpt)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}

	return &GCS{client}, nil
}

// Put creates a new cloud storage object or overwrites an existing one.
func (s *GCS) Put(ctx context.Context, bucket, objectName string, contents []byte, cacheable bool, contentType string) error {
	cacheControl := "public, max-age=86400"
	if !cacheable {
		cacheControl = "no-cache, max-age=0"
	}

	wc := s.client.Bucket(bucket).Object(objectName).NewWriter(ctx)
	wc.CacheControl = cacheControl
	if contentType != "" {
		wc.ContentType = contentType
	}

	if _, err := wc.Write(contents); err != nil {
		return fmt.Errorf("storage.Writer.Write: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("storage.Writer.Close: %w", err)
	}

	return nil
}

// Put creates a new cloud storage object or overwrites an existing one from dir.
func (s *GCS) FPut(ctx context.Context, parent, name, filePath string, cacheAble bool, contentType string) error {
	cacheControl := "public, max-age=86400"
	if !cacheAble {
		cacheControl = "no-cache, max-age=0"
	}

	wc := s.client.Bucket(parent).Object(name).NewWriter(ctx)
	wc.CacheControl = cacheControl
	if contentType != "" {
		wc.ContentType = contentType
	}

	// Open local file.
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %w", err)
	}

	if _, err := wc.Write(contents); err != nil {
		return fmt.Errorf("storage.Writer.Write: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("storage.Writer.Close: %w", err)
	}

	return nil
}

// Delete deletes a cloud storage object, returns nil if the object was
// successfully deleted, or of the object doesn't exist.
func (s *GCS) Delete(ctx context.Context, bucket, objectName string) error {
	if err := s.client.Bucket(bucket).Object(objectName).Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			// Object doesn't exist; presumably already deleted.
			return nil
		}
		return fmt.Errorf("storage.DeleteObject: %w", err)
	}
	return nil
}

// Get returns the contents for the given object. If the object does not
// exist, it returns ErrNotFound.
func (s *GCS) Get(ctx context.Context, bucket, object string) ([]byte, error) {
	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: conf.GCS.AcecssID,
		PrivateKey:     []byte(conf.GCS.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(1 * time.Minute),
	})

	if err != nil {
		return nil, fmt.Errorf("storage.SignedURL: %w", err)
	}

	return []byte(url), nil
}

func (s *GCS) ReSignedURLWithReplace(ctx context.Context, parent, object string) (string, error) {
	// Get the object metadata to check if it has expired.
	obj := s.client.Bucket(parent).Object(object)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		log.Printf("failed to get object attributes: %v\n", err)
		return "", nil // skip if the object doesn't exist (maybe it was deleted)
	}

	// log.Printf("obj info: %+v\n", attrs)
	// log.Printf("obj updated: %v > %v? %v\n", time.Since(attrs.Updated.UTC()), expiry, time.Since(attrs.Updated) > expiry)
	// log.Printf("time now: %v\n", time.Now())
	// log.Printf("one minute later: %v\n", time.Now().Add(1*time.Minute))
	// log.Printf("one minute later in UTC: %v\n", time.Now().In(time.UTC).Add(1*time.Minute))

	if time.Since(attrs.Updated) > TestDuration {
		// The URL has expired, generate a new signed URL with the same expiration time.
		expirationTime := time.Now().Add(TestDuration) // Set expiration time to 1 minutes from now
		newURL, err := storage.SignedURL(parent, obj.ObjectName(), &storage.SignedURLOptions{
			GoogleAccessID: conf.GCS.AcecssID,
			PrivateKey:     []byte(conf.GCS.PrivateKey),
			Method:         "GET",
			Expires:        expirationTime,
		})
		if err != nil {
			return "", fmt.Errorf("failed to generate signed URL: %v", err)
		}

		attrsToUpdate := storage.ObjectAttrsToUpdate{
			CustomTime: expirationTime.Local(),
		}
		_, err = obj.Update(ctx, attrsToUpdate)
		if err != nil {
			return "", fmt.Errorf("failed to update object attributes: %v", err)
		}

		return newURL, nil
	}

	return "", nil
}

func (s *GCS) ReSignedURL(ctx context.Context, parent, object, existingUrl string) (string, error) {
	// handle if existing url is empty
	if existingUrl == "" {
		url, err := storage.SignedURL(parent, object, &storage.SignedURLOptions{
			GoogleAccessID: conf.GCS.AcecssID,
			PrivateKey:     []byte(conf.GCS.PrivateKey),
			Method:         "GET",
			Expires:        time.Now().Add(1 * time.Minute),
		})

		if err != nil {
			log.Printf("storage.SignedURL: %v\n", err)
			return "", fmt.Errorf("storage.SignedURL: %w", err)
		}

		return url, nil
	}

	pinfo := PresigneUrlInfoGCS{}

	// validate obj expiry
	// parse last modified to time
	// exampleURl: http://127.0.0.1:9000/dci-auth/sample.jpeg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LINh4iBxsPqfwXkf%2F20230618%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230618T151930Z&X-Amz-Expires=20&X-Amz-SignedHeaders=host&X-Amz-Signature=7ced0793a4cf6ea2494cb5e11bba9c8a096081a3daf35f14076a1ed88eb6efb8
	err := utils.ParseUrl(existingUrl, &pinfo)
	if err != nil {
		log.Printf("failed to parse url : %v\n", err)
		return "", err
	}

	parseTime, err := utils.ParseUnix(pinfo.Expires)
	if err != nil {
		log.Printf("failed to parse date : %v\n", err)
		return "", err
	}

	convertExpires, err := strconv.Atoi(pinfo.Expires)
	if err != nil {
		log.Printf("failed to convert expires : %v\n", err)
		return "", err
	}

	if time.Since(parseTime) > time.Duration(convertExpires) {
		url, err := storage.SignedURL(parent, object, &storage.SignedURLOptions{
			GoogleAccessID: conf.GCS.AcecssID,
			PrivateKey:     []byte(conf.GCS.PrivateKey),
			Method:         "GET",
			Expires:        time.Now().Add(1 * time.Minute),
		})

		if err != nil {
			log.Printf("storage.SignedURL: %v\n", err)
			return "", fmt.Errorf("storage.SignedURL: %w", err)
		}

		return url, nil
	}

	return existingUrl, nil
}
