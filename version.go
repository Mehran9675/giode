package giode

// AppVersion and AppBuild are the app's version and build number from
// giode.toml (app.version, app.build): the giode CLI sets them when it
// builds or runs a desktop app (with -ldflags -X). They're empty in a
// plain go build.
var (
	AppVersion string
	AppBuild   string
)
