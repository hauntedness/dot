package types

import "flag"

var componentFlags = flag.NewFlagSet("component", flag.PanicOnError)

func init() {
	componentFlags.String("name", "", "name of the component")
	componentFlags.String("kind", "struct", "kind of the component")
}

var providerFlags = flag.NewFlagSet("provider", flag.PanicOnError)

func init() {
	providerFlags.String("name", "", "name of the provider")
	providerFlags.String("param", "", "param settings of the provider")
}
