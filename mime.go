package misc

import (
	"mime"
	"strings"
)

// builtinTypesLower 内置的常见 MIME 类型到默认扩展名映射（全部小写）。
var builtinTypesLower = map[string][]string{
	// https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/img
	// https://developer.mozilla.org/zh-CN/docs/Web/Media/Formats/Image_types
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	"image/apng":               {".apng"},
	"image/avif":               {".avif"},
	"image/gif":                {".gif"},
	"image/jpeg":               {".jpeg", ".jpg", ".jfif", ".pjpeg", ".pjp"},
	"image/png":                {".png"},
	"image/svg+xml":            {".svg"},
	"image/webp":               {".webp"},
	"image/bmp":                {".bmp"},
	"image/x-icon":             {".ico", ".cur"},
	"image/vnd.microsoft.icon": {".ico"},
	"image/tiff":               {".tif", ".tiff"},

	// https://developer.mozilla.org/zh-CN/docs/Web/Media/Formats/Containers#browser_compatibility
	"audio/3gpp":   {".3gp"},
	"audio/aac":    {".aac"},
	"audio/flac":   {".flac"},
	"audio/mpeg":   {".mpg", ".mpeg"},
	"audio/mp3":    {".mp3"},
	"audio/mp4":    {".mp4", ".m4a"},
	"audio/ogg":    {".oga", ".ogg"},
	"audio/wav":    {".wav"},
	"audio/webm":   {".webm", ".weba"},
	"audio/midi":   {".mid", ".midi"},
	"audio/x-midi": {".mid", ".midi"},
	"audio/opus":   {".opus"},

	// https://developer.mozilla.org/zh-CN/docs/Web/Media/Formats/Containers#browser_compatibility
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	"video/3gpp":      {".3gp"},
	"video/mpeg":      {".mpeg", ".mpg"},
	"video/mp4":       {".mp4", ".m4v", ".m4p"},
	"video/ogg":       {".ogv", ".ogg"},
	"video/quicktime": {".mov", ".qt"},
	"video/webm":      {".webm"},
	"video/x-msvideo": {".avi"},
	//"video/mp2t":      ".ts",
}

// ExtensionByType 根据 MIME 类型返回一个匹配的扩展名（包含点，如 .png）。
// 优先使用内置表，其次回退到标准库 mime.ExtensionsByType。
// 若无法识别，返回空字符串。
func ExtensionByType(mimeType string) string {
	extensions, ok := builtinTypesLower[mimeType]
	if ok {
		return extensions[0]
	}
	extensions, err := mime.ExtensionsByType(mimeType)
	if err == nil && len(extensions) > 0 {
		return extensions[0]
	}
	return ""
}

// TypeByExtension 根据扩展名（可不含点）返回对应的 MIME 类型。
// 优先匹配内置表，不存在时回退到标准库 mime.TypeByExtension。
func TypeByExtension(ext string) string {
	if ext == "" {
		return ""
	}
	if ext[0] != '.' {
		ext = "." + ext
	}
	for typ, exts := range builtinTypesLower {
		for _, ext2 := range exts {
			if ext2 == ext {
				return typ
			}
		}
	}
	return mime.TypeByExtension(ext)
}

// CharsetByType 针对 text/* 类型返回 "charset=utf-8"，否则返回空字符串。
func CharsetByType(typ string) string {
	if strings.HasPrefix(typ, "text/") {
		return "charset=utf-8"
	}
	return ""
}
