package handler

import "github.com/lkeix/vis-gode-dependency/testdata/usecase"

type hogeHandler struct {
	hogeUsecase usecase.HogeUsecase
}

func NewHogeHandler(hogeUsecase usecase.HogeUsecase) *hogeHandler {
	return &hogeHandler{
		hogeUsecase: hogeUsecase,
	}
}

func (h *hogeHandler) Hoge() {
	h.hogeUsecase.Hoge()
}
