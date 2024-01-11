package routeutils

import "github.com/diegoclair/go_boilerplate/infra/logger"

type Utils interface {
	Resp() ResponseUtils
	Req() RequestUtils
}

type routeUtils struct {
	resp ResponseUtils
	req  RequestUtils
}

func New(log logger.Logger) Utils {
	return &routeUtils{
		resp: newResponseUtils(log),
		req:  newRequestUtils(),
	}
}

func (r *routeUtils) Resp() ResponseUtils {
	return r.resp
}

func (r *routeUtils) Req() RequestUtils {
	return r.req
}
