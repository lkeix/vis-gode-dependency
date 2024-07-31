package infrastructure

type Hoge interface {
	Hoge()
}

type hoge struct{}

var _ Hoge = &hoge{}

func NewHoge() *hoge {
	return &hoge{}
}

func (h *hoge) Hoge() {

}
