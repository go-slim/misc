// Package misc 提供通用的实用工具函数与类型：
//  - 摘要/密码：MD5、Sha1、Sha256、PasswordHash/Verify
//  - 函数组合：Call/CallWith、Wrap/WrapWith、MustCall/MustCallWith
//  - MIME 解析：ExtensionByType、TypeByExtension、CharsetByType
//  - 模板插值：Strtr/Tmpl、Interpolate、TagFunc
//  - 堆栈信息：Stack（包含源码行）
//  - 零拷贝转换：UnsafeBytesToString、UnsafeStringToBytes 等（需谨慎使用）
//  - 零值判断：IsZero（支持指针递归判断）
package misc
