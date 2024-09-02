package images

// CreateResult stores result data about uploaded and created image of method [Service.Create]
type CreateResult struct {
	Slug   string `json:"slug"`
	Mime   string `json:"mime"`
	Ext    string `json:"ext"`
	Size   int64  `json:"size"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func (r *CreateResult) FromModel(img Image) {
	r.Slug = img.Slug
	r.Mime = img.Mime
	r.Ext = img.Ext
	r.Size = img.Size
	r.Width = img.Width
	r.Height = img.Height
	r.URL = img.Url
}
