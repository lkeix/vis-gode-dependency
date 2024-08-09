package infrastructure

import "github.com/lkeix/vis-gode-dependency/testdata/domain/repository"

type hoge struct{}

var _ repository.Hoge = &hoge{}

func NewHoge() *hoge {
	return &hoge{}
}

func (h *hoge) Hoge() {

}
