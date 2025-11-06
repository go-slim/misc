// Package misc provides common utility functions and types for Go applications.
//
// # Digest and Password Hashing
//
// Cryptographic hash functions for data integrity and password management:
//   - MD5: Fast hash for checksums (not cryptographically secure)
//   - Sha1: SHA-1 hash algorithm
//   - Sha256: SHA-256 hash algorithm (recommended for security)
//   - PasswordHash: Secure password hashing using bcrypt
//   - PasswordVerify: Verify password against bcrypt hash
//
// # Function Composition
//
// Utilities for composing and chaining functions with error handling:
//   - Call: Execute multiple functions sequentially, stop on first error
//   - CallG: Execute multiple functions with shared parameter, stop on first error
//   - Wrap: Combine functions into a single reusable function
//   - WrapG: Combine functions with parameters into a single reusable function
//
// # MIME Type Utilities
//
// MIME type and file extension mapping (deprecated, consider using github.com/h2non/filetype):
//   - ExtensionByType: Get file extension from MIME type
//   - TypeByExtension: Get MIME type from file extension
//   - CharsetByType: Get charset for text/* MIME types (limited functionality)
//
// # Template Interpolation
//
// Fast string template processing with placeholder substitution:
//   - Substitute: Simple {key} placeholder replacement with map data
//   - Interpolate: Custom delimiter placeholder replacement with map data
//   - Tmpl: Low-level template engine with custom TagFunc for advanced use cases
//   - TagFunc: Function type for custom tag processing
//
// # Generic Utilities
//
// Type-safe generic functions for common operations:
//   - Zero: Get zero value for any type
//   - Ptr: Create pointer to value (useful for literals)
//   - Nil: Get typed nil pointer
//   - IsZero: Check if value is zero (supports deep checking)
//   - IsNil: Check if value is nil (works with interfaces and reference types)
//   - Coalesce: Return first non-zero value from arguments
//
// # Math Utilities
//
// Generic math functions for ordered types:
//   - MinMax: Return minimum and maximum of two values
//   - Clamp: Constrain value to range [min, max]
//
// # Zero-Copy Conversions
//
// Unsafe zero-copy conversions between string and []byte (use with extreme caution):
//   - BytesToString: Convert []byte to string without allocation
//   - StringToBytes: Convert string to []byte without allocation
//
// Warning: These functions use unsafe operations. The converted values share
// underlying memory. Modifying the source after conversion causes undefined behavior.
//
// # Stack Tracing (Deprecated)
//
// Stack trace utilities with source code lines (deprecated, use runtime.Stack or debug.Stack):
//   - Stack: Get stack trace with source code lines
//
// # Example Usage
//
//	// Password hashing
//	hash, _ := misc.PasswordHash("mypassword")
//	ok := misc.PasswordVerify("mypassword", hash)
//
//	// Function composition
//	err := misc.Call(
//		func() error { fmt.Println("step 1"); return nil },
//		func() error { fmt.Println("step 2"); return nil },
//	)
//
//	// Template interpolation
//	result, _ := misc.Substitute("Hello {name}!", map[string]any{"name": "World"})
//
//	// Generic utilities
//	value := misc.Coalesce("", "default", "fallback") // returns "default"
//	min, max := misc.MinMax(5, 3) // returns (3, 5)
package misc
