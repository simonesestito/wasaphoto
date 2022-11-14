package photo

type EntityPhoto struct {
	Id          []byte `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	AuthorId    []byte `json:"authorId"`
	PublishDate string `json:"publishDate"`
}
