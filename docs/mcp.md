# MCP server

Giode can run a [Model Context Protocol](https://modelcontextprotocol.io) server inside your
app. AI agents connected to it can see the UI and use it the way a person would: read what's
on screen, click, type, take screenshots. That's useful for automated testing, scripted
demos, or letting an assistant work with your app.

```go
app := giode.New(giode.Config{
	Title: "My app",
	MCP:   true,
})
```

The server listens on `http://127.0.0.1:7600/mcp`. Connect a client to it, for example
Claude Code:

```sh
claude mcp add --transport http my-app http://127.0.0.1:7600/mcp
```

To turn the server on in any Giode app without changing code, set the environment
variable `GIODE_MCP=1`. `GIODE_MCP_ADDR` changes the address.

## Built-in tools

| Tool | What it does |
| --- | --- |
| `ui_snapshot` | Lists what's on screen: text, buttons, fields and other controls, each with an id, label, value, position and the actions it accepts. Optional `kind` and `query` filters. |
| `ui_click` | Clicks a button, tab or menu item, or toggles a checkbox. |
| `ui_type` | Replaces a text field's content, optionally pressing Enter afterwards. |
| `ui_set_value` | Sets a checkbox (`true`/`false`), slider (0 to 1) or select (option text or index). |
| `screenshot` | Returns the window as a PNG image. |
| `state_list`, `state_set` | Read and change [named states](state.md#named-state). |
| `set_theme` | Switches between light and dark mode. |
| `window` | Minimize, maximize, restore, fullscreen, center or raise the window. |

Actions go through the same code as real input: click handlers run and the UI updates.

## Adding tools

Write a tool as a Go function. The fields of its input struct become the tool's parameters:

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

- A parameter is required unless its `json` tag has `omitempty`.
- `description:"..."` and `enum:"a,b,c"` tags describe a parameter to the agent.
- Return a string, any value that encodes to JSON, or an `*mcp.Result` (for example
  `mcp.Image(png, "image/png")`). Returning an error reports a failed tool call.

Tools run on the server's goroutines, not the UI goroutine. Change UI data through `State`
(safe from any goroutine), `app.Post(fn)` or `app.RunOnUI(ctx, fn)`; see
[Threading](threading.md).

## Making controls visible to agents

Built-in components appear in `ui_snapshot` automatically. Controls without visible text
need a name: `IconButton(...).Label("Delete")`, `Input.Text(...).Label("Email")`,
`Pressable(...).Label("Open")`.

A component of your own appears when it calls `elements.RecordSemantic` in its `Layout`, and
the tools can act on it when it has the matching method: `Click()`, `ChangeText(string)`,
`Toggle()` with `Checked()`, `Change(float32)`, or `Options()` with `Choose(int)`.

## Options

```go
giode.Config{
	MCP: true,
	MCPOptions: giode.MCPOptions{
		Addr:     "127.0.0.1:7600", // listen address
		Token:    "secret",         // require "Authorization: Bearer secret"
		Stdio:    false,            // use stdin/stdout instead, for clients that start the app
		ReadOnly: false,            // only the tools that look, none that act
		// DisableBuiltinTools, Name, Version, Instructions, AllowedOrigins, Timeout
	},
}
```

## Security

The server only listens on the loopback interface by default, so other machines can't
reach it. It rejects requests from web pages on other origins, so a website open in your
browser can't talk to it. If you bind it to another interface, set `Token`.

The `mcp` package doesn't depend on the rest of Giode and can serve tools on its own, over
HTTP (`Server.ServeHTTP`, `ListenAndServe`) or stdio (`ServeStdio`).
