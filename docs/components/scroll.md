# Scroll

A vertical scroll container. Its content can be any height; scroll with the mouse wheel,
touch, or by dragging the scrollbar.

```go
feed := giode.Scroll(func() giode.Element {
	return giode.Stack(giode.Styles{Gap: 8}, posts()...)
}, giode.Styles{ScrollBar: giode.ScrollBar{ThumbColor: "theme.primary"}})

// in the view, give it a height:
giode.Box(giode.Styles{Height: 300}, feed)
```

The content function runs every frame. For long lists of equal-height rows, use
[VirtualList](virtuallist.md), which only builds the visible rows.

## Methods

| Method | Description |
| --- | --- |
| `Scroll(content, st...) *Scroll` | Creates the container. |
| `OffsetDp()`, `ScrollToDp(v)` | Scroll position in dp. |
| `Offset()`, `ScrollTo(v)` | Scroll position in pixels. |
| `Progress() float32` | 0 at the top, 1 at the bottom; 0 when everything fits. |
| `AtBottom() bool` | At the end, or everything fits. |
| `VisibleDp(top, height) bool` | Whether a range of the content, in dp, is at least partly on screen. |
| `Visible(top, height) bool` | The same in pixels. |
| `ContentHeight()`, `ViewportHeight()` | Measured on the last frame, in pixels. |
| `OnReachBottom(fn) *Scroll` | Called once each time the end is reached. |

## Loading more at the end

```go
feed.OnReachBottom(func() {
	go func() {
		more := fetchPage(next)
		app.Post(func() { posts = append(posts, more...) })
	}()
})
```

`OnReachBottom` also fires when the content doesn't fill the view, which loads the first
pages until it does.

## Skipping off-screen work

With rows of known height, skip building the expensive ones that aren't visible:

```go
feed := giode.Scroll(func() giode.Element {
	rows := make([]giode.Element, len(items))
	for i := range items {
		top := float32(i * 80)
		if !feed.VisibleDp(top, 80) {
			rows[i] = giode.Box(giode.Styles{Height: 80}) // placeholder
			continue
		}
		rows[i] = card(items[i])
	}
	return giode.Stack(giode.Styles{}, rows...)
})
```

Position checks use the previous frame's scroll offset, so a row becomes real one frame
after it scrolls into view.

## Styles

| Property | Role |
| --- | --- |
| `Background` | Fill behind the content. |
| `ScrollBar` | `Width`, `TrackColor`, `ThumbColor`, `Radius`, `MinThumbLength` (dp). |
