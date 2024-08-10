package infrastructure

import "github.com/lkeix/vis-gode-dependency/testdata/usecase/queryservice"

type piyo struct {
}

var _ queryservice.Piyo = &piyo{}

func NewPiyo() *piyo {
	return &piyo{}
}

func (p *piyo) Piyo() {

}
