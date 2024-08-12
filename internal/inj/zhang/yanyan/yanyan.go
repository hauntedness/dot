package yanyan

//go:generate dot-ioc
var _ int

//go:ioc component
type YanYan struct{}

//go:ioc --name yanYan
func NewYanYan() YanYan {
	return YanYan{}
}
