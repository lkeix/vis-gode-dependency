package usecase

import "github.com/lkeix/vis-gode-dependency/testdata/infrastructure"

type HogeUsecase interface {
	Hoge()
}

type hogeUsecase struct {
	hogeInfra infrastructure.Hoge
}

var _ HogeUsecase = &hogeUsecase{}

func NewHogeUsecase(hi infrastructure.Hoge) *hogeUsecase {
	return &hogeUsecase{
		hogeInfra: hi,
	}
}

func (h *hogeUsecase) Hoge() {
	h.hogeInfra.Hoge()
}
