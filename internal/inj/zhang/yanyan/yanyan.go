package yanyan

//go:generate digen -cmd=scan
var _ int

//go:ioc component
type YanYan struct{}

//go:ioc provider
func NewYanYan() YanYan {
	return YanYan{}
}
