package types

type DNodeState int64

const (
	IsInjected DNodeState = iota << 1
	IsValue
	IsPointer
	IsCyclicReferred
	IsInvalid
)

type Component struct {
	// State represent
	State DNodeState
}

// Identical 检查两个Components是否一样
func (*Component) Identical(other *Component) bool {
	panic("Not Implemented")
}

// DNode represent 1 dependency tree node.
type DNode struct {
	// 看起来需要两层索引, 一层放所有的provider
	// 一层放provider需要的Components
	Provider Provider
	// Dependencies of current node
	Dependencies []*DNode
	// Pos is the position of current node in DTree.Values
	Pos int
}

// IsLeaf denoted that current node has no dependencies
func (d *DNode) IsLeaf() bool {
	return len(d.Dependencies) == 0
}

// DTree
type DTree struct {
	Root *DNode
	// Components stores all nodes
	Components []*Component
}

func (tree *DTree) Component(node *DNode) *Component {
	return tree.Components[node.Pos]
}

func (tree *DTree) AppendObjects(components ...*Component) {
	add := []*Component{}
LoopInputs:
	for _, c := range components {
		for _, dep := range tree.Components {
			if c.Identical(dep) {
				continue LoopInputs
			}
		}
		add = append(add, c)
	}
	tree.Components = append(tree.Components, add...)
}

func (tree *DTree) WriteNode(node *DNode) string {
	// 获取component node的类型
	// 调用各自的实现
	// 抛出
	panic("Not Implemented")
}
