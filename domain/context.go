package domain

import "context"

type ContextDomain interface {
	Context() context.Context
	IsNew() bool
	NeedUpdate() bool
	updated()
}

var _ ContextDomain = (*baseContextDomain)(nil)

func newDomain(ctx context.Context) baseContextDomain {
	return baseContextDomain{
		ctx:   ctx,
		isNew: true,
	}
}

func loadDomain(ctx context.Context) baseContextDomain {
	return baseContextDomain{
		ctx: ctx,
	}
}

type baseContextDomain struct {
	ctx        context.Context
	isNew      bool
	needUpdate bool
}

func (b *baseContextDomain) updated() {
	b.needUpdate = true
}

func (b *baseContextDomain) Context() context.Context {
	return b.ctx
}

func (b *baseContextDomain) IsNew() bool {
	return b.isNew
}

func (b *baseContextDomain) NeedUpdate() bool {
	return b.needUpdate
}
