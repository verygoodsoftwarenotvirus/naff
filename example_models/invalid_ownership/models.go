package cyclic

type A struct {
	Name string

	_META_ uintptr `belongs_to:"C"`
}
