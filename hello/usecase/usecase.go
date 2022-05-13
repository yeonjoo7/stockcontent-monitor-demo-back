package usecase

import (
	"context"
	"github.com/google/uuid"
	"stockcontent-monitor-demo-back/domain"
	"time"
)

var _ domain.HelloUseCase = (*helloUseCase)(nil)

func NewHelloUseCase(
	helloRepo domain.HelloRepository,
	timeout time.Duration,
) domain.HelloUseCase {
	return &helloUseCase{
		helloRepo: helloRepo,
		dbTimeout: timeout,
	}
}

type helloUseCase struct {
	helloRepo domain.HelloRepository
	dbTimeout time.Duration
}

func (h *helloUseCase) GenerateHello(ctx context.Context, value string) (saved domain.HelloDomain, err error) {
	c, cancel := h.timeoutContext(ctx)
	defer cancel()

	saved = domain.NewHello(c, value)
	err = h.helloRepo.Save(saved)
	if err != nil {
		saved = nil
	}
	return
}

func (h *helloUseCase) GetHello(ctx context.Context, id uuid.UUID) (res domain.HelloDomain, err error) {
	c, cancel := h.timeoutContext(ctx)
	defer cancel()

	res = h.helloRepo.GetOne(c, id)
	if res == nil {
		err = domain.ErrItemNotFound
	}
	return
}

func (h *helloUseCase) timeoutContext(ctx context.Context) (c context.Context, cancel func()) {
	c, cancel = context.WithTimeout(ctx, h.dbTimeout)
	return
}
