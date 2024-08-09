package usecase

import (
	"github.com/lkeix/vis-gode-dependency/testdata/domain/repository"
)

type HogeUsecase interface {
	Hoge()
}

type hogeUsecase struct {
	hogeInfra repository.Hoge
}

var _ HogeUsecase = &hogeUsecase{}

func NewHogeUsecase(hi repository.Hoge) *hogeUsecase {
	return &hogeUsecase{
		hogeInfra: hi,
	}
}

func (h *hogeUsecase) Hoge() {
	h.hogeInfra.Hoge()
}
