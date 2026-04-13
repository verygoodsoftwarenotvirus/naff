package templates

import "embed"

// FS contains all embedded template files.
//
//go:embed domain/*.tmpl domain_keys/*.tmpl domain_fakes/*.tmpl domain_converters/*.tmpl domain_mock/*.tmpl manager/*.tmpl repository/*.tmpl codegen/*.tmpl migrations/*.tmpl proto/*.tmpl grpc/*.tmpl authorization/*.tmpl sqlc/*.tmpl build/*.tmpl project/*.tmpl project/scripts/*.tmpl ios/project/*.tmpl ios/models/*.tmpl ios/api/*.tmpl ios/auth/*.tmpl frontend/shared/*.tmpl frontend/consumer/project/*.tmpl frontend/consumer/layout/*.tmpl frontend/consumer/domain/*.tmpl frontend/consumer/entity/*.tmpl frontend/admin/project/*.tmpl frontend/admin/layout/*.tmpl frontend/admin/domain/*.tmpl frontend/admin/entity/*.tmpl
var FS embed.FS
