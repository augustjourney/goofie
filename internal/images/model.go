package images

import (
	"api/pkg/config"
	"fmt"
	"gorm.io/gorm"
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
	return i.WithPath().WithURL(cfg.CdnURL)
}

// WithMetadata sets metadata to Image instance and returns it.
func (i *Image) WithMetadata(metadata Metadata) *Image {
	i.Width = metadata.Width
	i.Height = metadata.Height
	return i
}
