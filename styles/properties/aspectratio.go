package properties

// AspectRatio fixes the element's width/height ratio (CSS
// aspect-ratio). When one axis has a resolved size and the other hugs
// content, the missing axis is derived from the ratio.
type AspectRatio float32
