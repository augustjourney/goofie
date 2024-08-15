package images

import "github.com/google/uuid"

type UploadImageDTO struct {
	Path     string
	Name     string
	Mime     string
	Ext      string
	Size     int64
	AuthorID int
}

func (d *UploadImageDTO) ToModel() *Image {
	return &Image{
		Ext:      d.Ext,
		Mime:     d.Mime,
		Size:     d.Size,
		Name:     d.Name,
		Slug:     uuid.New().String(),
		AuthorID: uint(d.AuthorID),
	}
}
