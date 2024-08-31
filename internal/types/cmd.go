package types

import "flag"

var componentFlags = flag.NewFlagSet("component", flag.PanicOnError)

func init() {
	componentFlags.String("name", "", "name of the component")
	componentFlags.String("kind", "struct", "kind of the component")
}
