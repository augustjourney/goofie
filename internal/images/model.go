package images

import (
	"api/pkg/config"
	"api/pkg/files"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
)

// Image stores data about image and represents table «images» in db
type Image struct {
	gorm.Model
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Mime     string `json:"mime"`
	Ext      string `json:"ext"`
	Path     string `json:"path"`
	Bucket   string `json:"bucket"`
	Provider string `json:"provider"`
	Url      string `json:"url"`
	Size     int64  `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	AuthorID uint   `json:"author_id"`
}

// FromFileHeader sets default values to image from received file header.
func (i *Image) FromFileHeader(ctx context.Context, file *multipart.FileHeader) *Image {
	i.Ext = files.GetExtension(ctx, file.Filename)
	i.Mime = file.Header.Get("Content-Type")
	i.Size = file.Size
	i.Name = file.Filename

	return i
}

// WithAuthor sets Author ID to Image instance and returns it.
func (i *Image) WithAuthor(authorID uint) *Image {
	i.AuthorID = authorID
	return i
}

// GetFilename returns image filename that is build from Slug and Ext
func (i *Image) GetFilename() string {
	return fmt.Sprintf("%s%s", i.Slug, i.Ext)
}

// WithPath sets Path to Image instance and returns it.
func (i *Image) WithPath() *Image {
	i.Path = fmt.Sprintf("/%s/%s", i.Bucket, i.GetFilename())
	return i
}

// WithURL sets Url to Image instance and returns it.
func (i *Image) WithURL(cdnURL string) *Image {
	i.Url = fmt.Sprintf("%s%s", cdnURL, i.Path)
	return i
}

// WithDefaults sets default fields such as bucket and provider to Image instance and returns it.
func (i *Image) WithDefaults(cfg *config.Config) *Image {
	i.Bucket = cfg.S3Bucket
	i.Provider = cfg.S3Provider
	i.Slug = uuid.NewString()
	return i.WithPath().WithURL(cfg.CdnURL)
}

// WithMetadata sets metadata to Image instance and returns it.
func (i *Image) WithMetadata(metadata Metadata) *Image {
	i.Width = metadata.Width
	i.Height = metadata.Height
	return i
}
