package photo

import "github.com/simonesestito/wasaphoto/service/api/route"

type Controller struct {
	Service Service
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{}
}
