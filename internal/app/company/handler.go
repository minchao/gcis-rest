package company

import (
	"context"
	"net/http"

	"github.com/minchao/gcis-rest/internal/pkg/gcisclient"
	"github.com/minchao/gcis-rest/internal/pkg/restutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/form"
	"github.com/minchao/go-gcis/gcis"
	"gopkg.in/go-playground/validator.v9"
)

var (
	chiAdapter *chiadapter.ChiLambda
	client     = gcisclient.New()
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if chiAdapter == nil {
		r := chi.NewRouter()
		r.Get("/companies", getCompaniesByKeyword)
		r.Get("/companies/{id:[0-9]+}", getCompany)

		chiAdapter = chiadapter.New(r)
	}

	return chiAdapter.Proxy(req)
}

type company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *company) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getCompany(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	info, _, err := client.Company.GetBasicInformation(context.Background(),
		&gcis.CompanyBasicInformationInput{BusinessAccountingNO: id})
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, &restutil.ErrorResponse{Message: "Unexpected error"})

		return
	}
	if info == nil {
		render.Status(r, http.StatusNotFound)
		_ = render.Render(w, r, &restutil.ErrorResponse{Message: "Company not found"})

		return
	}

	_ = render.Render(w, r, &company{
		ID:   info.BusinessAccountingNO,
		Name: info.CompanyName,
	})
}

type companiesResponse struct {
	Data []*company `json:"data"`
}

func (c *companiesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type companiesByKeywordOptions struct {
	Keyword string `form:"keyword" validate:"required"`
	Status  string `form:"status"`
	Limit   int    `form:"limit" validate:"gt=0"`
	Offset  int    `form:"offset"`
}

func (opt *companiesByKeywordOptions) isValid() error {
	return validator.New().Struct(opt)
}

func getCompaniesByKeyword(w http.ResponseWriter, r *http.Request) {
	var opt companiesByKeywordOptions
	err := form.NewDecoder().Decode(&opt, r.URL.Query())
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &restutil.ErrorResponse{Message: http.StatusText(http.StatusBadRequest)})
		return
	}
	if err := opt.isValid(); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		_ = render.Render(w, r, &restutil.ErrorResponse{Message: http.StatusText(http.StatusUnprocessableEntity)})
		return
	}
	if opt.Status == "" {
		opt.Status = "01"
	}

	list := make([]*company, 0)

	out, _, err := client.Company.SearchByKeyword(context.Background(),
		&gcis.CompanyByKeywordInput{
			CompanyName:   opt.Keyword,
			CompanyStatus: opt.Status,
			SearchOptions: gcis.SearchOptions{
				Top:  opt.Limit,
				Skip: opt.Offset,
			},
		})
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, &restutil.ErrorResponse{Message: "Unexpected error"})
		return
	}
	for _, info := range out {
		c := &company{
			ID:   info.BusinessAccountingNO,
			Name: info.CompanyName,
		}
		list = append(list, c)
	}

	_ = render.Render(w, r, &companiesResponse{Data: list})
}
