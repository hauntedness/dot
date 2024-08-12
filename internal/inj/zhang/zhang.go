package zhang

import "github.com/hauntedness/dot/internal/inj/zhang/yanyan"

//go:generate dot-ioc
var _ int

//go:ioc component
type Zhang struct{}

//go:ioc
func NewZhang(yanyan yanyan.YanYan) Zhang {
	return Zhang{}
}
