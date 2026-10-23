# MCP server

Giode can run a [Model Context Protocol](https://modelcontextprotocol.io) server inside your
app, so AI agents can see and drive it while it runs — for testing, automation, or giving an
assistant hands on your app. Turn it on with one option:

```go
app := giode.New(giode.Config{
	Title: "My app",
	MCP:   true,
})
```

The server listens on `http://127.0.0.1:7600/mcp` (Streamable HTTP). Point any MCP client
at it, for example Claude Code:

```bash
claude mcp add --transport http my-app http://127.0.0.1:7600/mcp
```

Setting the environment variable `GIODE_MCP=1` enables it in any Giode app without a code
change; `GIODE_MCP_ADDR` overrides the address.

## Built-in tools

| Tool | What it does |
| --- | --- |
| `ui_snapshot` | The visible UI as nodes: id, kind, label, value, bounds, accepted actions. Filter by `kind` or `query`. |
| `ui_click` | Clicks a button, tab or menu item, or toggles a checkbox. |
| `ui_type` | Replaces an input's text, optionally submitting it. |
| `ui_set_value` | Sets a checkbox (`true`/`false`), slider (0–1) or select (option or index). |
| `screenshot` | The window as a PNG. |
| `state_list` / `state_set` | Reads and sets named states (see below). |
| `set_theme` | Switches light/dark mode. |
| `window` | Minimize, maximize, fullscreen, center, raise, or `resize` to a `width` x `height` in dp - `{"action": "resize", "width": 400, "height": 840}` tests a phone layout on a desktop. |

Actions run on the UI thread, exactly as if the user did them: handlers fire and the UI
updates.

## Your own tools

Describe a tool with a typed Go function; its input struct becomes the JSON schema:

```go
app.MCPTool(mcp.Func("add_todo", "Adds an item to the todo list.",
	func(ctx context.Context, in struct {
		Title string `json:"title" description:"The todo's title"`
		Done  bool   `json:"done,omitempty"`
	}) (any, error) {
		todos.Update(func(t []Todo) []Todo { return append(t, Todo{in.Title, in.Done}) })
		return "added", nil
	}))
```

Fields are required unless tagged `omitempty`; `description` and `enum:"a,b"` tags document
them. Return a string, any JSON-encodable value, or an `*mcp.Result` (e.g. `mcp.Image`).
Handlers run on server goroutines: change the UI through `State` (safe from any goroutine)
or `app.RunOnUI(ctx, fn)`.

## Named state

States created with `giode.NamedState` are visible to `state_list` and settable through
`state_set`:

```go
var todos = giode.NamedState("todos", []Todo{})
```

## Labels for automation

Text-less controls should get a name: `IconButton(...).Label("Delete")`,
`Input.Text(...).Label("Email")`. Custom components join the snapshot with
`elements.RecordSemantic`, and become actionable by implementing `Click()`,
`ChangeText(string)`, `Toggle()`/`Checked()`, `Change(float32)` or
`Options()`/`Choose(int)`.

## Options and security

```go
MCPOptions: giode.MCPOptions{
	Addr:     "127.0.0.1:7600", // listen address
	Token:    "secret",         // require "Authorization: Bearer secret"
	Stdio:    false,            // serve stdin/stdout instead, for clients that launch the app
	ReadOnly: false,            // only snapshot/screenshot/state_list
	DisableBuiltinTools: false, // only your own tools
},
```

The server binds to loopback by default and rejects browser requests from non-localhost
origins (DNS-rebinding protection). Set `Token` before binding to any other interface. The
`mcp` package is independent of the UI and can serve tools on its own.
