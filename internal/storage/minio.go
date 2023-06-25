package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/vldcreation/sample-cron-go/internal/utils"
)

var (
	// Compile-time check to verify implements interface.
	_            Storage = (*Minio)(nil)
	cacheControl         = "public, max-age=86400"
)

type Minio struct {
	client *minio.Client
}

func NewMinio() (Storage, error) {
	client, err := minio.New(conf.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Minio.AccessKey, conf.Minio.SecretKey, ""),
		Secure: conf.Minio.UseSSL,
		Region: conf.Minio.Region,
	})
	if err != nil {
		log.Fatalf("error occured while initialize minio %v", err.Error())

		return nil, err
	}
	return &Minio{
		client: client,
	}, nil
}

func (m *Minio) Get(ctx context.Context, bucket, object string) ([]byte, error) {
	_, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	// check object storage exist
	// return CustomError if object not found
	_, err = m.client.StatObject(ctx, bucket, object, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("failed to get object : %v\n", err)
		return nil, nil // skip error if object not found
	}

	// default expiration time is 7 days
	// if you want to change expiration time, you can change it in time.Hour*24*7
	// log.Printf("test time : %v\n", time.Duration(time.Now().Local().AddDate(0, 0, 7).Day()*int(time.Millisecond)))
	urlObject, err := m.client.PresignedGetObject(ctx, bucket, object, Test10Seconds, nil)
	if err != nil {
		return nil, err
	}

	return []byte(urlObject.Redacted()), nil
}

func (m *Minio) Put(ctx context.Context, bucket, object string, data []byte, cacheAble bool, contentType string) error {
	reader := bytes.NewReader(data)

	if !cacheAble {
		cacheControl = "no-cache, max-age=0"
	}

	// do validation to make sure bucket and object name is valid
	if err := s3utils.CheckValidBucketName(bucket); err != nil {
		return err
	}

	if err := s3utils.CheckValidObjectName(object); err != nil {
		return err
	}

	// check available bucket
	// bucket available or not || create bucket
	ok, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	// create bucket if not available
	// bucket region should be same with minio region
	if !ok {
		bucketName := conf.Minio.Bucket

		if err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: conf.Minio.Region,
		}); err != nil {
			return err
		}
	}

	_, err = m.client.PutObject(ctx, bucket, object, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType:  contentType,
		CacheControl: cacheControl,
	})

	return err
}

func (m *Minio) FPut(ctx context.Context, bucket, object, filePath string, cacheAble bool, contentType string) error {
	if !cacheAble {
		cacheControl = "no-cache, max-age=0"
	}

	// do validation to make sure bucket and object name is valid
	if err := s3utils.CheckValidBucketName(bucket); err != nil {
		return err
	}

	if err := s3utils.CheckValidObjectName(object); err != nil {
		return err
	}

	// check available bucket
	// bucket available or not || create bucket
	ok, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	// create bucket if not available
	// bucket region should be same with minio region
	if !ok {
		bucketName := conf.Minio.Bucket

		if err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: conf.Minio.Region,
		}); err != nil {
			return err
		}
	}

	_, err = m.client.FPutObject(ctx, bucket, object, filePath, minio.PutObjectOptions{
		ContentType:  contentType,
		CacheControl: cacheControl,
	})

	return err
}

func (m *Minio) Delete(ctx context.Context, bucket, object string) error {
	_, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	// check object storage exist
	// return CustomError if object not found
	_, err = m.client.StatObject(ctx, bucket, object, minio.StatObjectOptions{})
	if err != nil {
		return nil // skip error if object not found
	}
	return m.client.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
}

func (m *Minio) ReSignedURLWithReplace(ctx context.Context, parent, object string) (string, error) {
	// check available bucket
	// more info : https://docs.min.io/docs/golang-client-api-reference.html#BucketExists
	obj, err := m.client.StatObject(ctx, parent, object, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("failed to get object : %v\n", err)
		return "", nil // skip error if object not found maybe it's already deleted
	}

	// validate obj expiry
	if time.Since(obj.LastModified) > Test10Seconds {
		// generated new Url
		// expirationTime := time.Now().Add(1 * time.Minute) // Set expiration time to 1 minutes from now
		newUrlObject, err := m.client.PresignedGetObject(ctx, parent, object, Test20Seconds, nil)
		if err != nil {
			return "", err
		}

		// use copy object to itself to update last-modified value strategy
		_, err = m.client.CopyObject(ctx, minio.CopyDestOptions{
			Bucket: parent,
			Object: object,
			UserMetadata: map[string]string{
				"Content-Type":  obj.ContentType,
				"Cache-Control": obj.Metadata.Get("Cache-Control"),
			},
			ReplaceMetadata: true,
		}, minio.CopySrcOptions{
			Bucket:               parent,
			Object:               object,
			VersionID:            obj.VersionID,
			MatchUnmodifiedSince: obj.LastModified,
			MatchETag:            obj.ETag,
		})

		if err != nil {
			return "", fmt.Errorf("failed to update object attributes: %v", err)
		}

		return newUrlObject.String(), nil
	}

	return "", nil

}

func (m *Minio) ReSignedURL(ctx context.Context, parent, object, existingUrl string) (string, error) {
	// check available bucket
	// more info : https://docs.min.io/docs/golang-client-api-reference.html#BucketExists
	_, err := m.client.StatObject(ctx, parent, object, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("failed to get object : %v\n", err)
		return "", nil // skip error if object not found maybe it's already deleted
	}

	if existingUrl != "" {
		newUrlObject, err := m.client.PresignedGetObject(ctx, parent, object, Test20Seconds, nil)
		if err != nil {
			return "", err
		}

		return newUrlObject.Redacted(), nil
	}

	pinfo := PresignUrlInfoS3{}

	// validate obj expiry
	// parse last modified to time
	// exampleURl: http://127.0.0.1:9000/dci-auth/sample.jpeg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LINh4iBxsPqfwXkf%2F20230618%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230618T151930Z&X-Amz-Expires=20&X-Amz-SignedHeaders=host&X-Amz-Signature=7ced0793a4cf6ea2494cb5e11bba9c8a096081a3daf35f14076a1ed88eb6efb8
	err = utils.ParseUrl(existingUrl, &pinfo)
	if err != nil {
		log.Printf("failed to parse url : %v\n", err)
		return "", err
	}

	parseTime, err := utils.ParseDate(pinfo.X_AMZ_DATE)
	if err != nil {
		log.Printf("failed to parse date : %v\n", err)
		return "", err
	}

	convertExpires, err := strconv.Atoi(pinfo.X_AMZ_EXPIRES)
	if err != nil {
		log.Printf("failed to convert expires : %v\n", err)
		return "", err
	}

	if time.Since(parseTime) > time.Duration(convertExpires) {
		// generated new Url
		// expirationTime := time.Now().Add(1 * time.Minute) // Set expiration time to 1 minutes from now
		newUrlObject, err := m.client.PresignedGetObject(ctx, parent, object, Test20Seconds, nil)
		if err != nil {
			return "", err
		}

		return newUrlObject.Redacted(), nil
	}

	return existingUrl, nil

}
