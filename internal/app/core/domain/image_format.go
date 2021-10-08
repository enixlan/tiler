package domain

type ImageFormat string

const (
	ImageFormatPNG = "png"
)

func (format ImageFormat) String() string {
	return string(format)
}

func (format ImageFormat) ContentType() (string, bool) {
	switch format {
	case ImageFormatPNG:
		return "image/png", true
	default:
		return "", false
	}
}
