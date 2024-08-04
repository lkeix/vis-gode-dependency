package handler

import "github.com/lkeix/vis-gode-dependency/testdata/usecase"

type fugaHandler struct {
	fugaUsecase usecase.FugaUsecase
}

func NewFugaHandler(fugaUsecase usecase.FugaUsecase) *fugaHandler {
	return &fugaHandler{
		fugaUsecase: fugaUsecase,
	}
}

func (f *fugaHandler) Fuga() {
	f.fugaUsecase.Fuga()
}
