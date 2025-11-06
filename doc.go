// Package misc provides common utility functions and types:
//   - Digest/Password: MD5, Sha1, Sha256, PasswordHash/Verify
//   - Function composition: Call/CallG, Wrap/WrapG
//   - MIME parsing: ExtensionByType, TypeByExtension, CharsetByType
//   - Template interpolation: Strtr/Tmpl, Interpolate, TagFunc
//   - Stack information: Stack (including source code lines) [Deprecated]
//   - Zero-copy conversion: UnsafeBytesToString, UnsafeStringToBytes, etc. (use with caution)
//   - Zero value judgment: IsZero (supports pointer recursive judgment)
package misc
