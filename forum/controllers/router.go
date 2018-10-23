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
	//==================
	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/:path1",
		models.Creator,
	},
	//=======================
	Route{
		"ForumGetOne",
		strings.ToUpper("Get"),
		"/api/forum/:slug/details",
		models.ForumGetOne,
	},

	Route{
		"ForumGetThreads",
		strings.ToUpper("Get"),
		"/api/forum/:slug/threads",
		models.ForumGetThreads,
	},

	Route{
		"ForumGetUsers",
		strings.ToUpper("Get"),
		"/api/forum/:slug/users",
		models.ForumGetUsers,
	},

	Route{
		"PostGetOne",
		strings.ToUpper("Get"),
		"/api/post/:id/details",
		models.PostGetOne,
	},

	Route{
		"PostUpdate",
		strings.ToUpper("Post"),
		"/api/post/:id/details",
		models.PostUpdate,
	},

	Route{
		"PostsCreate",
		strings.ToUpper("Post"),
		"/api/thread/:slug_or_id/create",
		models.PostsCreate,
	},

	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		models.Status,
	},
	//=========================
	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/:path1/:path2",
		models.Creator,
	},
	//=======================
	Route{
		"ThreadGetOne",
		strings.ToUpper("Get"),
		"/api/thread/:slug_or_id/details",
		models.ThreadGetOne,
	},

	Route{
		"ThreadGetPosts",
		strings.ToUpper("Get"),
		"/api/thread/:slug_or_id/posts",
		models.ThreadGetPosts,
	},

	Route{
		"ThreadUpdate",
		strings.ToUpper("Post"),
		"/api/thread/:slug_or_id/details",
		models.ThreadUpdate,
	},

	Route{
		"ThreadVote",
		strings.ToUpper("Post"),
		"/api/thread/:slug_or_id/vote",
		models.ThreadVote,
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
