package db

type Series int

const (
	Unspecified Series = iota
	Original
	TheFallOfMaxPayne
)

func (s Series) String() string {
	switch s {
	case Original:
		return "1"
	case TheFallOfMaxPayne:
		return "The Fall of Max Payne"
	default:
		return "0"
	}
}

type Excerpt struct {
	Series  int    `json:"series"`
	Part    string `json:"part"`
	Chapter string `json:"chapter"`
	Excerpt string `json:"excerpt"`
}
