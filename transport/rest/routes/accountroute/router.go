package accountroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
)

const RouteName = "accounts"

const (
	rootRoute          = ""
	accountByID        = "/:account_uuid/"
	accountBalanceByID = "/:account_uuid/balance"
)

type UserRouter struct {
	ctrl      *Handler
	routeName string
}

func NewRouter(ctrl *Handler, routeName string) *UserRouter {
	return &UserRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *UserRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(r.routeName)

	router.POST(rootRoute, r.ctrl.handleAddAccount).
		Summary("Add a new account").
		Read(viewmodel.AddAccount{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}})

	router.POST(accountBalanceByID, r.ctrl.handleAddBalance).
		Summary("Add balance to an account").
		Description("Add balance to an account by account_uuid").
		Read(viewmodel.AddBalance{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}})

	router.GET(rootRoute, r.ctrl.handleGetAccounts).
		Summary("Get all accounts").
		Description("Get all accounts with paginated response").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PaginatedResult[[]viewmodel.AccountResponse]{},
			},
		}).
		QueryParam("page", "number of page you want", goswag.StringType, false).
		QueryParam("quantity", "quantity of items per page", goswag.StringType, false)

	router.GET(accountByID, r.ctrl.handleGetAccountByID).
		Summary("Get account by ID").
		Description("Get account by it UUID value").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.AccountResponse{},
			},
		})
}
