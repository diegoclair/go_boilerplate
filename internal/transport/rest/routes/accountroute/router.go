package accountroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/viewmodel"
	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
)

const GroupRouteName = "accounts"

const (
	RootRoute               = ""
	AccountByIDRoute        = "/:account_uuid/"
	AccountBalanceByIDRoute = "/:account_uuid/balance"
)

type AccountRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *AccountRouter {
	return &AccountRouter{
		ctrl: ctrl,
	}
}

func (r *AccountRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(GroupRouteName)

	router.POST(RootRoute, r.ctrl.handleAddAccount).
		Summary("Add a new account").
		Read(viewmodel.AddAccount{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}})

	router.POST(AccountBalanceByIDRoute, r.ctrl.handleAddBalance).
		Summary("Add balance to an account").
		Description("Add balance to an account by account_uuid").
		Read(viewmodel.AddBalance{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}}).
		PathParam("account_uuid", "account uuid", goswag.StringType, true)

	router.GET(RootRoute, r.ctrl.handleGetAccounts).
		Summary("Get all accounts").
		Description("Get all accounts with paginated response").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PaginatedResponse[[]viewmodel.AccountResponse]{},
			},
		}).
		QueryParam("page", "number of page you want", goswag.StringType, false).
		QueryParam("quantity", "quantity of items per page", goswag.StringType, false)

	router.GET(AccountByIDRoute, r.ctrl.handleGetAccountByID).
		Summary("Get account by ID").
		Description("Get account by it UUID value").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.AccountResponse{},
			},
		}).
		PathParam("account_uuid", "account uuid", goswag.StringType, true)
}
