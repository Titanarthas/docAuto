package models

type Spec struct {
	ApiSpecs ApiSpecs
}

type ApiSpecs []ApiSpec

func (c ApiSpecs) Len() int           { return len(c) }
func (c ApiSpecs) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ApiSpecs) Less(i, j int) bool { return c[i].Path > c[j].Path }
