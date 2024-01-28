package routeutils

type routeUtils struct {
	resp ResponseUtils
	req  RequestUtils
}

func New() Utils {
	return &routeUtils{
		resp: newResponseUtils(),
		req:  newRequestUtils(),
	}
}

func (r *routeUtils) Resp() ResponseUtils {
	return r.resp
}

func (r *routeUtils) Req() RequestUtils {
	return r.req
}
