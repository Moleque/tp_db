package controllers

import (
	"strings"
	"tp_db/forum/models"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.RedirectTrailingSlash = true
	for _, route := range routes {
		var handler httprouter.Handle
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.Handle(route.Method, route.Pattern, handler)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		Index,
	},

	Route{
		"Clear",
		strings.ToUpper("Post"),
		"/api/service/clear",
		Clear,
	},

	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/create",
		models.ForumCreate,
	},

	Route{
		"ForumGetOne",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/details",
		ForumGetOne,
	},

	Route{
		"ForumGetThreads",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/threads",
		ForumGetThreads,
	},

	Route{
		"ForumGetUsers",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/users",
		ForumGetUsers,
	},

	Route{
		"PostGetOne",
		strings.ToUpper("Get"),
		"/api/post/{id}/details",
		PostGetOne,
	},

	Route{
		"PostUpdate",
		strings.ToUpper("Post"),
		"/api/post/{id}/details",
		PostUpdate,
	},

	Route{
		"PostsCreate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/create",
		PostsCreate,
	},

	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		Status,
	},

	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/{slug}/create",
		ThreadCreate,
	},

	Route{
		"ThreadGetOne",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/details",
		ThreadGetOne,
	},

	Route{
		"ThreadGetPosts",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/posts",
		ThreadGetPosts,
	},

	Route{
		"ThreadUpdate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/details",
		ThreadUpdate,
	},

	Route{
		"ThreadVote",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/vote",
		ThreadVote,
	},

	Route{
		"UserCreate",
		strings.ToUpper("Post"),
		"/api/user/:nickname/create",
		models.UserCreate,
	},

	Route{
		"UserGetOne",
		strings.ToUpper("Get"),
		"/api/user/:nickname/profile",
		models.UserGetOne,
	},

	Route{
		"UserUpdate",
		strings.ToUpper("Post"),
		"/api/user/:nickname/profile",
		models.UserUpdate,
	},
}
