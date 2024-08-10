package handler

import "github.com/lkeix/vis-gode-dependency/testdata/usecase/queryservice"

type piyoHandler struct {
	piyoqs queryservice.Piyo
}

func NewPiyoHandler(piyoqs queryservice.Piyo) *piyoHandler {
	return &piyoHandler{
		piyoqs: piyoqs,
	}
}

func (h *piyoHandler) Piyo() {
	h.piyoqs.Piyo()
}
