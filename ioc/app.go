package ioc

type app struct {
	*Container
}

var appInstance *app

func newApp(container *Container) *app {
	return &app{Container: container}
}

func App() *app {
	if appInstance == nil {
		appInstance = newApp(newContainer())
	}
	return appInstance
}
