package image

type WMType string

const (
	WM_TYPE_TEXT  WMType = `text`
	WM_TYPE_IMAGE WMType = `image`
)

// WatermarkOptions 水印选项
type WatermarkOptions struct {
	Watermark string `json:"watermark"` // 水印图片文件路径
	Type WMType `json:"type,omitempty"` // 水印类型
	Position Pos `json:"position"` // 水印的位置
	Padding int `json:"padding"` // 水印留的边白
	On bool `json:"on"` // 是否开启水印
}

func (w *WatermarkOptions) SetWatermark(watermark string, typ WMType) *WatermarkOptions {
	w.Watermark = watermark
	w.Type = typ
	return w
}

func (w *WatermarkOptions) SetPosition(position Pos) *WatermarkOptions {
	w.Position = position
	return w
}

func (w *WatermarkOptions) SetPadding(padding int) *WatermarkOptions {
	w.Padding = padding
	return w
}

func (w *WatermarkOptions) Enable() *WatermarkOptions {
	w.On = true
	return w
}

func (w *WatermarkOptions) IsEnabled() bool {
	return w.On && len(w.Watermark) > 0
}

func (w *WatermarkOptions) Disable() *WatermarkOptions {
	w.On = false
	return w
}

func (w *WatermarkOptions) SetOn(on bool) *WatermarkOptions {
	w.On = on
	return w
}

func (w *WatermarkOptions) CreateInstance() (*Watermark, error) {
	return NewWatermark(w.Watermark, w.Padding, w.Position)
}
