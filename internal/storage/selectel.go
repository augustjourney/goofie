package storage

import (
	"api/pkg/consts"
	"api/pkg/logger"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Selectel implements [S3] interface to interact with s3 selectel storage.
type Selectel struct {
	lastAuthTime time.Time
	authToken    string
	cfg          Config
}

// NewSelectel creates and returns a new [Selectel] instance.
func NewSelectel(cfg Config) *Selectel {
	return &Selectel{
		cfg: cfg,
	}
}

// Auth sends a request to selectel to get auth token for further upload requests.
// Auth token lives only 24 hours.
func (s *Selectel) Auth() error {
	// token lives 24 hours
	if time.Since(s.lastAuthTime).Hours() < 23 && s.authToken != "" {
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.selcdn.ru/auth/v1.0", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-User", s.cfg.AccessKeyID)
	req.Header.Set("X-Auth-Key", s.cfg.AccessKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	s.authToken = resp.Header.Get("x-auth-token")
	s.lastAuthTime = time.Now()
	return nil
}

// Upload uploads a file to selectel s3 storage.
func (s *Selectel) Upload(ctx context.Context, buffer io.ReadSeeker, bucket string, fileName string, mime string, expiry string) (string, error) {
	startUploadingAt := time.Now()

	err := s.Auth()
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	url := fmt.Sprintf("https://api.selcdn.ru/v1/SEL_207069/%s/%s", bucket, fileName)

	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Auth-Token", s.authToken)
	req.Header.Set("Content-Type", mime)

	if expiry != "" {
		req.Header.Set("X-Delete-At", expiry)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	uploadProcessingTime := time.Since(startUploadingAt).Milliseconds()
	ctx = context.WithValue(ctx, consts.UploadProcessingTimeKey, uploadProcessingTime)

	logger.Info(logger.Record{
		Message: fmt.Sprintf("upload to selectel success, bucket: %s, file: %s, expiry: %s", bucket, fileName, expiry),
		Context: ctx,
	})

	return buildFilePath(bucket, fileName), nil
}
