package zhang

import "github.com/hauntedness/dot/internal/inj/zhang/yanyan"

//go:generate digen -cmd=scan
var _ int

//go:ioc component
type Zhang struct{}

//go:ioc provider
func NewZhang(yanyan yanyan.YanYan) Zhang {
	return Zhang{}
}
