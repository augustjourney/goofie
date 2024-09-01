package images

// CreateResult stores result data about uploaded and created image of method [Service.Create]
type CreateResult struct {
	Slug   string `json:"slug"`
	Mime   string `json:"mime"`
	Ext    string `json:"ext"`
	Size   int64  `json:"size"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
