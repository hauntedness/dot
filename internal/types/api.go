package types

// Provider
type Provider interface {
	Provide() []*Component
	//
	Require() []*Component
	// 如何渲染
	Template() string
}
