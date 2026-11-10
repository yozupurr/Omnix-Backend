package profile

type template struct {
	attributes func(season int) map[string]any
	cosmetics  bool
}

var templates = map[string]template{
	"athena":   {attributes: athena, cosmetics: true},
	"profile0": {attributes: profile0},
}
