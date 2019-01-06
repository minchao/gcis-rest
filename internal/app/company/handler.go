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
	"github.com/minchao/go-gcis/gcis"
)

var (
	chiAdapter *chiadapter.ChiLambda
	client     = gcisclient.New()
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if chiAdapter == nil {
		r := chi.NewRouter()
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
