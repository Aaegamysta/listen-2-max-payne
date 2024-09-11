package db

type Series int

const (
	Unspecified Series = iota
	Original
	TheFallOfMaxPayne
)

type Excerpt struct {
	// ID string	maybe?
	Series  int    `json:"series"`
	Part    string `json:"part"`
	Chapter string `json:"chapter"`
	Excerpt string `json:"excerpt"`
}
