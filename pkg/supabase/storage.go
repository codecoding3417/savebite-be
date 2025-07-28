package supabase

import (
	"fmt"
	storagego "github.com/supabase-community/storage-go"
	"io"
	"savebite/internal/domain/env"
	"savebite/pkg/log"
)

type SupabaseItf interface {
	Upload(bucket, path string, contentType string, file io.Reader) (string, error)
}

type SupabaseStruct struct {
	client *storagego.Client
}

var Supabase = getSupabase()

func getSupabase() SupabaseItf {
	url := fmt.Sprintf("%s/storage/v1", env.AppEnv.SupabaseURL)
	client := storagego.NewClient(url, env.AppEnv.SupabaseSecret, nil)

	return &SupabaseStruct{client}
}

func (s *SupabaseStruct) Upload(bucket, path string, contentType string, file io.Reader) (string, error) {
	res, err := s.client.UploadFile(bucket, path, file, storagego.FileOptions{
		ContentType: &contentType,
	})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
			"res":   res,
		}, "[SupabaseStorage][Upload] Failed to upload file")
		return "", err
	}

	return s.client.GetPublicUrl(bucket, path).SignedURL, nil
}
