package domain

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/google/uuid"
	"stockcontent-monitor-demo-back/ent"
	"stockcontent-monitor-demo-back/util/entx"
)

// HelloData anemic model
type HelloData struct {
	Id    uuid.UUID
	Value string
}

type HelloDomain interface {
	ContextDomain
	Id() uuid.UUID
	Say() string
	SetValue(value string)
	RawData() HelloData
}

type HelloUseCase interface {
	GenerateHello(ctx context.Context, value string) (HelloDomain, error)
	GetHello(ctx context.Context, id uuid.UUID) (HelloDomain, error)
}

type HelloRepository interface {
	Save(domain HelloDomain) error
	GetOne(ctx context.Context, id uuid.UUID) HelloDomain
	Transaction(ctx context.Context, opts *sql.TxOptions, fn func(HelloTxRepository) error) error
}

type HelloTxRepository interface {
	entx.Tx
	HelloRepository
}

var _ HelloDomain = (*helloImpl)(nil)

func FromHelloEntity(ctx context.Context, entity *ent.Hello) HelloDomain {
	return &helloImpl{
		baseContextDomain: loadDomain(ctx),
		id:                entity.ID,
		value:             entity.Value,
	}
}

func NewHello(ctx context.Context, value string) HelloDomain {
	return &helloImpl{
		baseContextDomain: newDomain(ctx),
		id:                uuid.New(),
		value:             value,
	}
}

type helloImpl struct {
	baseContextDomain
	id    uuid.UUID
	value string
}

func (h *helloImpl) SetValue(value string) {
	defer h.updated()
	h.value = value
}

func (h *helloImpl) Id() uuid.UUID {
	return h.id
}

func (h *helloImpl) Say() string {
	return fmt.Sprintf("hello %s", h.value)
}

func (h *helloImpl) RawData() HelloData {
	return HelloData{
		Id:    h.id,
		Value: h.value,
	}
}
