package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Selectel struct {
	lastAuthTime time.Time
	authToken    string
	cfg          Config
}

func NewSelectel(cfg Config) *Selectel {
	return &Selectel{
		cfg: cfg,
	}
}

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

func (s *Selectel) Upload(ctx context.Context, buffer io.ReadSeeker, bucket string, fileName string, mime string, expiry string) (string, error) {
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

	fmt.Println("Uploaded file with status", resp.StatusCode)

	return buildFilePath(bucket, fileName), nil
}
