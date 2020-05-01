# naff

NAFF is a code generation tool. It stands for:
&emsp;**N**ot
&emsp;**A**nother
&emsp;**F**ucking
&emsp;**F**ramework

## selfishware

NAFF was built for me, by me. If NAFF is useful to you too, I consider that simultaneously cool and collateral. If it isn't helpful to you, I don't care. 

I wasn't even going to make this code public in any way, and was going to just provide binaries, but I don't run binaries on my machine whose code I'm not allowed to read, and I don't suggest you do so either.

I'm honestly not interested in feature ideas/pull requests/complaints/suggestions/etc. As far as I am personally concerned, NAFF has exactly one user, and I am them. Even this README is for my own self (save for this section). 

## state of the codebase

The codebase is rough, and will probably remain rough for a very long time. It started life as mostly generated, and has been mangled many hundreds of times by poorly constructed regular expressions. 

By the time this README is in master, the code will be generating everything appropriately. At some point I will go back and add unit tests for all the functions to test all the cases and only then will I consider refactoring the codebase.

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

Running `naff gen gitlab.com/verygoodsoftwarenotvirus/addressbook` will trigger a series of prompts from the CLI, and then it will generate code.
 
NAFF writes code only on clean slates, and as such will nuke the output directory from orbit. So if I output to `gitlab.com/verygoodsoftwarenotvirus/naff`, or some other precious path, I'm going to have a bad time.

The structs you defined can only have a small selection of types, basically anything that can be defined as a constant in Go. No type aliases, no embedded structs, no `time.Time`s, no slices. Pointers are allowed, though. Here are all the allowed types:

```
bool
int
int8
int16
int32
int64
uint
uint8
uint16
uint32
uint64
uintptr  // should only be used for the _META_ field
float32
float64
string
```

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

You can indicate that an object belongs to nobdy with the struct tag `belongs_to:"__nobody__"`

Note that something can belong to a user and also to another struct.