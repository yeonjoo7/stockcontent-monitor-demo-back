package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var templateFunc = template.FuncMap{
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
}

var domainModelTemplate = template.Must(template.New("domainModel").
	Funcs(templateFunc).
	Parse(`package domain

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"stockcontent-monitor-demo-back/ent"
	"stockcontent-monitor-demo-back/util/entx"
)

// {{.}}Data anemic model
type {{.}}Data struct {
}

type {{.}}Domain interface {
	ContextDomain
	RawData() {{.}}Data
}

type {{.}}UseCase interface {
}

type {{.}}Repository interface {
	Save(domain {{.}}Domain) error
	Transaction(ctx context.Context, opts *sql.TxOptions, fn func({{.}}TxRepository) error) error
}

type {{.}}TxRepository interface {
	entx.Tx
	{{.}}Repository
}

var _ {{.}}Domain = &{{. | lower}}Impl{}

func From{{.}}Entity(ctx context.Context, entity *ent.{{.}}) {{.}}Domain {
	return &{{. | lower}}Impl{
		baseContextDomain: loadDomain(ctx),
	}
}

func New{{.}}(ctx context.Context) {{.}}Domain {
	return &{{. | lower}}Impl{
		baseContextDomain: newDomain(ctx),
	}
}

type {{. | lower}}Impl struct {
	baseContextDomain
}

func (h *{{. | lower}}Impl) RawData() {{.}}Data {
	return {{.}}Data{
	}
}
`))

var handlerTemplate = template.Must(template.New("handler").
	Funcs(templateFunc).
	Parse(`package handler

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"net/http"
	"stockcontent-monitor-demo-back/domain"
)


const (
	tag = "{{. | upper}}-CONTROLLER"
)

var {{.}}ControllerProvider = wire.NewSet(
	wire.Struct(new({{.}}Controller), "*"),
)

type {{.}}Controller struct {
	UseCase domain.{{.}}UseCase
}

func (c *{{.}}Controller) Bind(e *echo.Echo) {
	// routing sample code
	e.GET("/sample.{{.}}", c.sample)
}

func (c *{{.}}Controller) sample(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "{{.}} sample routing",
	})
}
`))

var useCaseTemplate = template.Must(template.New("useCase").
	Funcs(templateFunc).
	Parse(`package usecase

import (
	"context"
	"stockcontent-monitor-demo-back/domain"
	"time"
)

var _ domain.{{.}}UseCase = &{{. | lower}}UseCase{}

func New{{.}}UseCase(
	{{. | lower}}Repo domain.{{.}}Repository,
	timeout time.Duration,
) domain.{{.}}UseCase {
	return &{{. | lower}}UseCase{
		{{. | lower}}Repo: {{. | lower}}Repo,
		dbTimeout: timeout,
	}
}

type {{. | lower}}UseCase struct {
	{{. | lower}}Repo domain.{{.}}Repository
	dbTimeout time.Duration
}

func (u *{{. | lower}}UseCase) timeoutContext(ctx context.Context) (c context.Context, cancel func()) {
	c, cancel = context.WithTimeout(ctx, u.dbTimeout)
	return
}
`))

var repositoryTemplate = template.Must(template.New("repository").
	Funcs(templateFunc).
	Parse(`package repository

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"stockcontent-monitor-demo-back/domain"
	"stockcontent-monitor-demo-back/ent"
)

var _ domain.{{.}}TxRepository = &mysqlRepo{}

func NewH{{.}}MySQLRepository(db *ent.Client) domain.{{.}}Repository {
	return &mysqlRepo{
		cli: db,
	}
}

type mysqlRepo struct {
	cli *ent.Client
	tx  *ent.Tx
}

func (m *mysqlRepo) {{. | lower}}Client() *ent.{{.}}Client {
	if m.tx != nil {
		return m.tx.{{.}}
	}

	if m.cli != nil {
		return m.cli.{{.}}
	}

	return nil
}

func (m *mysqlRepo) Transaction(ctx context.Context, opts *sql.TxOptions, fn func(domain.{{.}}TxRepository) error) error {
	tx, err := m.cli.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(new{{.}}TxRepository(tx)); err != nil {
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

func (m *mysqlRepo) GetTx() *ent.Tx {
	return m.tx
}

func (m *mysqlRepo) Save(domain domain.{{.}}Domain) (err error) {
	// uncommented if you define Id field
	// data := domain.RawData()

	if domain.IsNew() {
		_, err = m.{{. | lower}}Client().Create().
			// set field code
			// SetID(data.Id).SetValue(data.Value).
			Save(domain.Context())
	}

	if domain.NeedUpdate() {
		// need Id field todo define Id
		_, err = m.{{. | lower}}Client().UpdateOneID(0 /*data.Id*/).
			// set field code
			// SetValue(data.Value).
			Save(domain.Context())
	}

	return
}

func new{{.}}TxRepository(tx *ent.Tx) domain.{{.}}TxRepository {
	return &mysqlRepo{tx: tx}
}
`))

func main() {
	name := os.Args[1]
	lowerName := strings.ToLower(name)
	os.MkdirAll("domain", os.ModePerm)
	domainFile, err := os.Create(fmt.Sprintf("domain/%s.go", lowerName))
	if err != nil {
		panic(err)
	}

	handlerPath := fmt.Sprintf("%s/handler", lowerName)
	useCasePath := fmt.Sprintf("%s/usecase", lowerName)
	repositoryPath := fmt.Sprintf("%s/repository", lowerName)
	os.MkdirAll(handlerPath, os.ModePerm)
	os.MkdirAll(useCasePath, os.ModePerm)
	os.MkdirAll(repositoryPath, os.ModePerm)

	handlerFile, err := os.Create(fmt.Sprintf("%s/http.go", handlerPath))
	if err != nil {
		panic(err)
	}

	useCaseFile, err := os.Create(fmt.Sprintf("%s/usecase.go", useCasePath))
	if err != nil {
		panic(err)
	}

	repositoryFile, err := os.Create(fmt.Sprintf("%s/mysql_repo.go", repositoryPath))
	if err != nil {
		panic(err)
	}

	domainModelTemplate.Execute(domainFile, name)
	handlerTemplate.Execute(handlerFile, name)
	useCaseTemplate.Execute(useCaseFile, name)
	repositoryTemplate.Execute(repositoryFile, name)
}
