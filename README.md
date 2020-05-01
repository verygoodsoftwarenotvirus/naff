# NAFF

NAFF is a code generation tool. It stands for:

&emsp;**N**ot<br>
&emsp;**A**nother<br>
&emsp;**F**ucking<br>
&emsp;**F**ramework<br>

## selfishware

NAFF was built by myself, for myself. If NAFF is useful to you too, I consider that simultaneously cool and incidental. If it isn't helpful to you, then I guess I'm sorry? Please don't be mad.

I wasn't even going to make this code public (see the below section), I was just going to provide binaries, but I don't run binaries on my machine whose code I'm not allowed to read, and I don't suggest you do so either.

I'm honestly not interested in feature ideas/pull requests/complaints/suggestions/etc and have consequently turned those features off in Gitlab. As far as I am personally concerned, NAFF has exactly one user, and I am them. Even this README is for my own self (save for this section). 

## state of the codebase

The codebase is rough, and will probably remain rough for a very long time. It started life as mostly generated, and has been mangled many hundreds of times by poorly constructed regular expressions. 

By the time this README is in master, the code will be generating everything appropriately. At some point I will go back and add unit tests for all the functions to test all the cases and only then will I consider refactoring the codebase.

There are some unit tests for things I just had to verify, but they're not run in CI or anything (another thing on my todo list), so they're likely broken.

## motivation

The goal for NAFF is to generate well-tested, fleshed-out web server repositories, so you can focus on the actual business logic of whatever you're trying to build, because you're not writing boilerplate.

I'm gonna write a blog post about this eventually, but the long and short of it is: I have an irrational distaste towards so-called "batteries included" web frameworks, and ORMS (almost equally). 

So I wrote an example CRUD server in Go, and then I wrote this tool to generate server codebases that are similar in (personally determined) quality to the example CRUD server codebase.

NAFF doesn't just generate code, it also generates:

- the Makefile and all relevant targets
- A service client
- Unit tests
- Integration tests
- Load tests
- Frontend (browser-driven) tests
- Database migrations/querying code
- All dependency injection code
- Dockerfiles
- Docker-compose files
- Routing code
- CI scripts
- a very very basic frontend in Svelte
- Telemetry (tracing/logging/metrics collection)
- Prometheus and Grafana configuration files
- Linter configuration

NAFF has support for multiple database providers. Currently those are:
   
- PostgreSQL
- Sqlite3
- MariaDB

Each of these may be (de)activated to your liking.

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
 
NAFF writes code only on clean slates, and as such will nuke the output directory from orbit. So if I tell NAFF to write output to some precious path, I'm going to have a bad time.

The structs you defined can only have a small selection of types, basically anything that can be defined as a constant in Go. No type aliases, no embedded structs, no `time.Time`, no slices. Pointers are allowed, though. Here are all the allowed types:

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

After generating a codebase, it's a good idea to run `make gamut`, which will generate all the initial files you're going to need to start writing the real code.

## advanced usage (ownership)

NAFF looks for a small, inconsistent smattering of flags in types. Let's say we wanted users to be able to set a birth year for their friend on creation, but not to change that birth year ever. You could accomplish something like that with this:

```
type Friend struct {
    FirstName string
    LastName string
    YearOfBirth uint `naff:"!editable"`
    MonthOfBirth uint
    DayOfBirth uint
}
```

A similar flag `!creatable` is also supported.

NAFF is capable of generating code that accounts for "ownership" between structs. The example I've been using is an old-school web discussion board, where you have the hierarchy `Forums > Subforums > Threads > Posts`. That hierarchy can be represented to NAFF like so:

```
type Forum struct {
    Name string
}

type Subforum struct {
    Name string
    _META_ uintptr `belongs_to:"Forum"`
}

type Thread struct {
    Name string
    _META_ uintptr `belongs_to:"Subforum,User"`
}

type Post struct {
    Content string
    _META_ uintptr `belongs_to:"Post,User"`
}
```

Note the `_META_` field. NAFF won't generate anything around this field, but looks to the struct tags of that field for information about how to handle types.

By default, all objects belong to a user, so you only need to indicate user ownership in the `_META_` field if the struct also belongs to another struct. You can indicate that an object belongs to nobody (effectively an enumeration) with the struct tag `belongs_to:"__nobody__"`:


```
type ValidZipCode struct {
    Code string
    _META_ uintptr `belongs_to:"__nobody__"`
}
```

Note that `Thread` and `Post` both belong to users and other structs. In this example, ordinary users can query posts and threads despite their ownership, because that matches the real-world analog. You can restrict something to a user with another tag:

```
type Friend struct {
    FirstName string
    LastName string
    YearOfBirth uint `naff:"!editable"`
    MonthOfBirth uint
    DayOfBirth uintstring
   _META_ uintptr `belongs_to:"Post,User"`
}

type FriendPhoto struct {
    PhotoURL string
    _META_ uintptr `restricted_to_user:"true"`
}
```

This will make it so that the generated code will only return a given user's data to that user.