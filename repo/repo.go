package repo

type Images struct {
	ImageName  string    `json:"name"`
	ImangeTags []*string `json:"tags"`
}

type Repository struct {
	Repos []*string `json:"repositories"`
}
