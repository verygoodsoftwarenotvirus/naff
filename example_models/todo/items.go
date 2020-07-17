package todo

type Item struct {
	Name    string
	Details string

	_META_ uintptr `restricted_to_user:"true",search_enabled:"true"`
}
