package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"stockcontent-monitor-demo-back/domain"
	"stockcontent-monitor-demo-back/ent"
	"stockcontent-monitor-demo-back/ent/hello"
)

var _ domain.HelloTxRepository = &mysqlEntRepo{}

type mysqlEntRepo struct {
	cli *ent.Client
	tx  *ent.Tx
}

func (m *mysqlEntRepo) helloClient() *ent.HelloClient {
	if m.tx != nil {
		return m.tx.Hello
	}

	if m.cli != nil {
		return m.cli.Hello
	}

	return nil
}

func (m *mysqlEntRepo) Transaction(ctx context.Context, fn func(domain.HelloTxRepository) error) error {
	tx, err := m.cli.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(&mysqlEntRepo{tx: tx}); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%v, rolling back transaction: %v", err, rerr)
			//err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%v, committing transaction: %v", err, err)
		//return errors.Wrapf(err, "committing transaction: %v", err)
	}
	return nil
}

func (m *mysqlEntRepo) GetTx() *ent.Tx {
	return m.tx
}

func (m *mysqlEntRepo) GetOne(ctx context.Context, id uuid.UUID) domain.HelloDomain {
	h, err := m.helloClient().
		Query().
		Where(hello.ID(id)).
		First(ctx)

	if err != nil {
		//todo notfound error handling ??
	}

	return domain.FromHelloEntity(ctx, h)
}

func (m *mysqlEntRepo) Save(domain domain.HelloDomain) (err error) {
	data := domain.RawData()
	if domain.IsNew() {
		_, err = m.helloClient().Create().
			SetID(data.Id).SetValue(data.Value).Save(domain.Context())
	}

	if domain.NeedUpdate() {
		_, err = m.helloClient().UpdateOneID(data.Id).
			SetValue(data.Value).Save(domain.Context())
	}

	return
}
