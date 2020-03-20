//go:generate stringer -type=Code -linecomment
package anytool

// Code error code defined
type Code int

// code 错误码
const (
	CodeFeatureNotSupport Code = 1000 + iota // feature Not support
	CodeOperationFailed                      // operate failed
	CodeInvalidArguments                     // Invalid arguments
	CodeExistYet                             // exist yet
)
