package controllers

import (
	"net/http"
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
		// handler = Logger(handler, route.Name)
		router.Handle(route.Method, route.Pattern, handler)
	}
	return router
}

func Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	return
}

func RouterOfCreating(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	path1 := params.ByName("path1")
	if path1 == "create" {
		models.ForumCreate(w, r, params)
	} else {
		models.ThreadCreate(w, r, params)
	}
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		Index,
	},

	// ===============
	// === Service ===
	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		models.Status,
	},

	Route{
		"Clear",
		strings.ToUpper("Post"),
		"/api/service/clear",
		models.Clear,
	},

	// =============
	// === Forum ===
	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/:path1",
		RouterOfCreating,
	},

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

	// ============
	// === Post ===
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

	// ==============
	// === Thread ===
	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/:path1/:path2",
		RouterOfCreating,
	},

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

	// ============
	// === User ===
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
