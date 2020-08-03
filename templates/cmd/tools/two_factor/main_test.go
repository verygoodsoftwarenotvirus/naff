package two_factor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()
		assert.NoError(t, RenderPackage(proj))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/\0/\0/\0`

		assert.Error(t, RenderPackage(proj))
	})
}

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"bufio"
	"fmt"
	totp "github.com/pquerna/otp/totp"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	zero  = "  ___   & / _ \\  &| | | | &| |_| | & \\___/  "
	one   = "    _    &  /_ |   &   | |   &  _| |_  & |_____| "
	two   = " ____   &|___ \\  &  __) | & / __/  &|_____| "
	three = "_____   &|___ /  &  |_ \\  & ___) | &|____/  "
	four  = " _   _   &| | | |  &| |_| |_ &|___   _ &    |_|  "
	five  = " ____   &| ___|  &|___ \\  & ___) | &|____/  "
	six   = "  __    & / /_   &| '_ \\  &| (_) | & \\___/  "
	seven = " _____  &|___  | &   / /  &  / /   & /_/    "
	eight = "  ___   & ( o )  & /   \\  &|  O  | & \\___/  "
	nine  = "  ___   & /   \\  &| (_) | & \\__, | &   /_/  "
)

var (
	lastChange  time.Time
	currentCode string

	numbers = [10][5]string{
		limitSlice(strings.Split(zero, "&")),
		limitSlice(strings.Split(one, "&")),
		limitSlice(strings.Split(two, "&")),
		limitSlice(strings.Split(three, "&")),
		limitSlice(strings.Split(four, "&")),
		limitSlice(strings.Split(five, "&")),
		limitSlice(strings.Split(six, "&")),
		limitSlice(strings.Split(seven, "&")),
		limitSlice(strings.Split(eight, "&")),
		limitSlice(strings.Split(nine, "&")),
	}
)

func limitSlice(in []string) (out [5]string) {
	if len(in) != 5 {
		panic("wut")
	}
	for i := 0; i < 5; i++ {
		out[i] = in[i]
	}
	return
}

func mustnt(err error) {
	if err != nil {
		panic(err)
	}
}

func clearTheScreen() {
	fmt.Println("\x1b[2J")
	fmt.Printf("\x1b[0;0H")
}

func buildTheThing(token string) string {
	var out string
	for i := 0; i < 5; i++ {
		if i != 0 {
			out += "\n"
		}
		for _, x := range strings.Split(token, "") {
			y, err := strconv.Atoi(x)
			if err != nil {
				panic(err)
			}
			out += "  "
			out += numbers[y][i]
		}
	}

	timeLeft := (30*time.Second - time.Since(lastChange).Round(time.Second)).String()
	out += fmt.Sprintf("\n\n%s\n", timeLeft)

	return out
}

func doTheThing(secret string) {
	t := strings.ToUpper(secret)
	n := time.Now().UTC()
	code, err := totp.GenerateCode(t, n)
	mustnt(err)

	if code != currentCode {
		lastChange = time.Now()
		currentCode = code
	}

	if !totp.Validate(code, t) {
		panic("this shouldn't happen")
	}

	clearTheScreen()
	fmt.Println(buildTheThing(code))
}

func requestTOTPSecret() string {
	var (
		token string
		err   error
	)

	if len(os.Args) == 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("token: ")
		token, err = reader.ReadString('\n')
		mustnt(err)
	} else {
		token = os.Args[1]
	}

	return token
}

func main() {
	secret := requestTOTPSecret()
	clearTheScreen()
	doTheThing(secret)
	every := time.Tick(1 * time.Second)
	lastChange = time.Now()

	for range every {
		doTheThing(secret)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildConstDeclarations()

		expected := `
package example

import ()

const (
	zero  = "  ___   & / _ \\  &| | | | &| |_| | & \\___/  "
	one   = "    _    &  /_ |   &   | |   &  _| |_  & |_____| "
	two   = " ____   &|___ \\  &  __) | & / __/  &|_____| "
	three = "_____   &|___ /  &  |_ \\  & ___) | &|____/  "
	four  = " _   _   &| | | |  &| |_| |_ &|___   _ &    |_|  "
	five  = " ____   &| ___|  &|___ \\  & ___) | &|____/  "
	six   = "  __    & / /_   &| '_ \\  &| (_) | & \\___/  "
	seven = " _____  &|___  | &   / /  &  / /   & /_/    "
	eight = "  ___   & ( o )  & /   \\  &|  O  | & \\___/  "
	nine  = "  ___   & /   \\  &| (_) | & \\__, | &   /_/  "
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildVarDeclarations()

		expected := `
package example

import (
	"strings"
	"time"
)

var (
	lastChange  time.Time
	currentCode string

	numbers = [10][5]string{
		limitSlice(strings.Split(zero, "&")),
		limitSlice(strings.Split(one, "&")),
		limitSlice(strings.Split(two, "&")),
		limitSlice(strings.Split(three, "&")),
		limitSlice(strings.Split(four, "&")),
		limitSlice(strings.Split(five, "&")),
		limitSlice(strings.Split(six, "&")),
		limitSlice(strings.Split(seven, "&")),
		limitSlice(strings.Split(eight, "&")),
		limitSlice(strings.Split(nine, "&")),
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildLimitSlice(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildLimitSlice()

		expected := `
package example

import ()

func limitSlice(in []string) (out [5]string) {
	if len(in) != 5 {
		panic("wut")
	}
	for i := 0; i < 5; i++ {
		out[i] = in[i]
	}
	return
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMustnt(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMustnt()

		expected := `
package example

import ()

func mustnt(err error) {
	if err != nil {
		panic(err)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildClearTheScreen(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildClearTheScreen()

		expected := `
package example

import (
	"fmt"
)

func clearTheScreen() {
	fmt.Println("\x1b[2J")
	fmt.Printf("\x1b[0;0H")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildTheThing(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildBuildTheThing()

		expected := `
package example

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func buildTheThing(token string) string {
	var out string
	for i := 0; i < 5; i++ {
		if i != 0 {
			out += "\n"
		}
		for _, x := range strings.Split(token, "") {
			y, err := strconv.Atoi(x)
			if err != nil {
				panic(err)
			}
			out += "  "
			out += numbers[y][i]
		}
	}

	timeLeft := (30*time.Second - time.Since(lastChange).Round(time.Second)).String()
	out += fmt.Sprintf("\n\n%s\n", timeLeft)

	return out
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDoTheThing(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildDoTheThing()

		expected := `
package example

import (
	"fmt"
	totp "github.com/pquerna/otp/totp"
	"strings"
	"time"
)

func doTheThing(secret string) {
	t := strings.ToUpper(secret)
	n := time.Now().UTC()
	code, err := totp.GenerateCode(t, n)
	mustnt(err)

	if code != currentCode {
		lastChange = time.Now()
		currentCode = code
	}

	if !totp.Validate(code, t) {
		panic("this shouldn't happen")
	}

	clearTheScreen()
	fmt.Println(buildTheThing(code))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequestTOTPSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildRequestTOTPSecret()

		expected := `
package example

import (
	"bufio"
	"fmt"
	"os"
)

func requestTOTPSecret() string {
	var (
		token string
		err   error
	)

	if len(os.Args) == 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("token: ")
		token, err = reader.ReadString('\n')
		mustnt(err)
	} else {
		token = os.Args[1]
	}

	return token
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMain()

		expected := `
package example

import (
	"time"
)

func main() {
	secret := requestTOTPSecret()
	clearTheScreen()
	doTheThing(secret)
	every := time.Tick(1 * time.Second)
	lastChange = time.Now()

	for range every {
		doTheThing(secret)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
