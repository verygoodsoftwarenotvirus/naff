package a

type Item struct {
	Name    string `naff:"createable,editable"`
	Details string `naff:"createable,editable"`
}

type Person struct {
	Name      string  `naff:"createable,editable"`
	Age       *uint64 `naff:"createable"`
	Gender    float32 `naff:"createable,editable"`
	Birthdate uint64  `naff:"createable,editable"`
}
