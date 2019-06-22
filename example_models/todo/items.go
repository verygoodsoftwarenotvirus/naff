package todo

type Item struct {
	Name    string `naff:"createable,editable"`
	Details string `naff:"createable,editable"`
}
