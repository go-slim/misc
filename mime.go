package misc

import (
	"mime"
	"slices"
	"strings"
)

// builtinTypesLower built-in common MIME types to default extensions mapping (all lowercase).
var builtinTypesLower = map[string][]string{
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/img
	// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types
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

	// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Containers#browser_compatibility
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

	// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Containers#browser_compatibility
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

// ExtensionByType returns a matching extension (including dot, like .png) based on MIME type.
// Prioritize built-in table, then fallback to standard library mime.ExtensionsByType.
// Returns empty string if unrecognized.
//
// Deprecated: Please use https://github.com/h2non/filetype library instead.
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

// TypeByExtension returns corresponding MIME type based on extension (can be without dot).
// Prioritize built-in table matching, fallback to standard library mime.TypeByExtension if not found.
//
// Deprecated: Please use https://github.com/h2non/filetype library instead.
func TypeByExtension(ext string) string {
	if ext == "" {
		return ""
	}
	if ext[0] != '.' {
		ext = "." + ext
	}
	for typ, exts := range builtinTypesLower {
		if slices.Contains(exts, ext) {
			return typ
		}
	}
	return mime.TypeByExtension(ext)
}

// CharsetByType returns "charset=utf-8" for text/* types, otherwise returns empty string.
//
// Deprecated: No longer supported.
func CharsetByType(typ string) string {
	if strings.HasPrefix(typ, "text/") {
		return "charset=utf-8"
	}
	return ""
}
