package data 

type Series int

const (
	Unspecified Series = iota
	Original 
	TheFallOfMaxPayne
)

type Excerpt struct {
	// ID string	maybe?
	Series string
	Part string
	Chapter string
	Excerpt string
}