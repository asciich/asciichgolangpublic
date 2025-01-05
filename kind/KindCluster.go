package kind

import "github.com/asciich/asciichgolangpublic"

type KindCluster struct {
	name string
	kind Kind
}

func NewKindCluster() (g *KindCluster) {
	return new(KindCluster)
}

func (g *KindCluster) GetKind() (kind Kind, err error) {

	return g.kind, nil
}

func (g *KindCluster) GetName() (name string, err error) {
	if g.name == "" {
		return "", asciichgolangpublic.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *KindCluster) MustGetKind() (kind Kind) {
	kind, err := g.GetKind()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kind
}

func (g *KindCluster) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return name
}

func (g *KindCluster) MustSetKind(kind Kind) {
	err := g.SetKind(kind)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (g *KindCluster) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (g *KindCluster) SetKind(kind Kind) (err error) {
	g.kind = kind

	return nil
}

func (g *KindCluster) SetName(name string) (err error) {
	if name == "" {
		return asciichgolangpublic.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
