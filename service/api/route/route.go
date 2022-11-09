package route

type Controller interface {
	ListRoutes() []Route
}

// Route is the description of an endpoint / API route.
// A controller will return the list of its routes,
// instead of registering them directly
// to be independent of the underlying routing library.
type Route interface {
	GetMethod() string
	GetPath() string
	IsSecure() bool
}

type AnonymousRoute struct {
	Method  string
	Path    string
	Handler Handler
}

func (route AnonymousRoute) GetMethod() string { return route.Method }
func (route AnonymousRoute) GetPath() string   { return route.Path }
func (route AnonymousRoute) IsSecure() bool    { return false }

type SecureRoute struct {
	Method  string
	Path    string
	Handler SecureHandler
}

func (route SecureRoute) GetMethod() string { return route.Method }
func (route SecureRoute) GetPath() string   { return route.Path }
func (route SecureRoute) IsSecure() bool    { return true }
