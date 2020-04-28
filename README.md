# naff

NAFF is a code generation tool. It stands for:
&emsp;**N**ot
&emsp;**A**nother
&emsp;**F**ucking
&emsp;**F**ramework

## selfishware

Full disclosure: NAFF was built for me, by me. If NAFF is useful to you too, that's cool, but also I don't care. I wasn't even going to make this code public in any way, and was going to just provide binaries, but I don't run binaries on my machine whose code I'm not allowed to read, and I don't suggest you do so either. 

The codebase itself is rough, and will probably remain rough for a very long time. I'm honestly not interested in feature ideas/pull requests/complaints/suggestions/etc. As far as I am personally concerned, NAFF has exactly one user, and I am them. Even this README is for my own self (save for this section).

## usage

Start by defining any arbitrary Go package with some types in it. For instance, say I have a Go package `gitlab.com/verygoodsoftwarenotvirus/addressbook`, with a file `types.go` with the following:

```
type Friend struct {
    FirstName string
    LastName string
    YearOfBirth uint
    MonthOfBirth uint
    DayOfBirth uint
}
```

Running `naff gen gitlab.com/verygoodsoftwarenotvirus/addressbook` will trigger a series of prompts from the CLI (all of which can be replaced with flags), and then it will generate code. NAFF deletes any folders present at the destination site, so if I output to `gitlab.com/verygoodsoftwarenotvirus/naff`, I'm going to have a bad time.

The structs you defined can only have a small selection of types, basically anything that can be defined as a constant in Go. No type aliases, no embedded structs, no `time.Time`s. Pointers are allowed, though.

## advanced usage

NAFF is also capable of generating code that accounts for "ownership" between structs. Take our earlier example, and expand upon it by allowing a user to attach photos to their friends. 

```
type Friend struct {
    FirstName string
    LastName string
    YearOfBirth uint
    MonthOfBirth uint
    DayOfBirth uint
}

type FriendPhoto struct {
    PhotoURL string
	_META_ uintptr `belongs_to:"Friend"`
}
```

Note the `_META_` field. NAFF won't generate anything around this field, but looks to the struct tags of that field for information about how to handle types.

For the above example, NAFF will generate code that requires some relevant `Friend` data in order to manipulate `FriendPhoto`s. By default, all objects belong to a user. 

You can indicate that an object belongs to nobody with the struct tag `belongs_to:"__nobody__"`