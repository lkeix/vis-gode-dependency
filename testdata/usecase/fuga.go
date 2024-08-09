package usecase

import (
	"github.com/lkeix/vis-gode-dependency/testdata/domain/repository"
)

type FugaUsecase interface {
	Fuga()
}

type fugaUsecase struct {
	hogeInfra repository.Hoge
}

var _ FugaUsecase = &fugaUsecase{}

func NewfugaUsecase(hi repository.Hoge) *fugaUsecase {
	return &fugaUsecase{
		hogeInfra: hi,
	}
}

func (h *fugaUsecase) Fuga() {
	h.hogeInfra.Hoge()
}
