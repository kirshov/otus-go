package internalhttp

import "net/http"

type route struct {
	URL     string
	Handler http.HandlerFunc
}

type routes []route

func getRoutes(app Application) routes {
	handlers := NewHandlers(app)

	return routes{
		route{
			URL:     "/",
			Handler: handlers.Home,
		},
		route{
			URL:     "/create",
			Handler: handlers.Create,
		},
		route{
			URL:     "/update",
			Handler: handlers.Update,
		},

		route{
			URL:     "/delete",
			Handler: handlers.Delete,
		},
		route{
			URL:     "/list",
			Handler: handlers.List,
		},
	}
}
