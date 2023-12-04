package providers

var providerMap = make(map[string]*provider)

type provider struct {
	Name        string
	bind        *Bind
	unbind      *UnBind
	refRequired bool
}

func Register(name string, bind Bind, unbind UnBind, refRequired bool) {
	providerMap[name] = &provider{
		Name:        name,
		bind:        &bind,
		unbind:      &unbind,
		refRequired: refRequired,
	}
}
