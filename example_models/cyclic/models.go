package cyclic

type A struct {
	Name string

	_META_ uintptr `belongs_to:"C"`
}

type B struct {
	Name string

	_META_ uintptr `belongs_to:"A"`
}

type C struct {
	Name string

	_META_ uintptr `belongs_to:"B"`
}
