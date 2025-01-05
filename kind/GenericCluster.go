package kind

import "github.com/asciich/asciichgolangpublic"

type GenericCluster struct {
	name string
	kind Kind
}

func NewGenericCluster() (g *GenericCluster) {
	return new(GenericCluster)
}

func (g *GenericCluster) GetKind() (kind Kind, err error) {

	return g.kind, nil
}

func (g *GenericCluster) GetName() (name string, err error) {
	if g.name == "" {
		return "", asciichgolangpublic.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GenericCluster) MustGetKind() (kind Kind) {
	kind, err := g.GetKind()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kind
}

func (g *GenericCluster) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return name
}

func (g *GenericCluster) MustSetKind(kind Kind) {
	err := g.SetKind(kind)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (g *GenericCluster) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (g *GenericCluster) SetKind(kind Kind) (err error) {
	g.kind = kind

	return nil
}

func (g *GenericCluster) SetName(name string) (err error) {
	if name == "" {
		return asciichgolangpublic.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
