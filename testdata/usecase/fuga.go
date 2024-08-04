package usecase

import "github.com/lkeix/vis-gode-dependency/testdata/infrastructure"

type FugaUsecase interface {
	Fuga()
}

type fugaUsecase struct {
	hogeInfra infrastructure.Hoge
}

var _ FugaUsecase = &fugaUsecase{}

func NewfugaUsecase(hi infrastructure.Hoge) *fugaUsecase {
	return &fugaUsecase{
		hogeInfra: hi,
	}
}

func (h *fugaUsecase) Fuga() {
	h.hogeInfra.Hoge()
}
