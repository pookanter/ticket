package ticket

import (
	"ticket/api/ticket/boards"
	"ticket/pkg/apikit"
)

func Router(api *apikit.API) {
	boards.Router(api)
}
