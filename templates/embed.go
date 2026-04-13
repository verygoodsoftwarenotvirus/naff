package templates

import "embed"

// FS contains all embedded template files.
//
//go:embed domain/*.tmpl domain_keys/*.tmpl domain_fakes/*.tmpl domain_converters/*.tmpl domain_mock/*.tmpl manager/*.tmpl repository/*.tmpl codegen/*.tmpl migrations/*.tmpl proto/*.tmpl grpc/*.tmpl authorization/*.tmpl sqlc/*.tmpl build/*.tmpl
var FS embed.FS
