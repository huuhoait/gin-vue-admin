package router

// Router is the plugin-local singleton wired into initialize.Router.
var Router = struct {
	SkyAgent SkyAgentRouter
}{}
