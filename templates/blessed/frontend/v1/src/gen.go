package frontendsrc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, projectName wordsmith.SuperPalabra, types []models.DataType) error {
	files := map[string]func() []byte{
		"frontend/v1/src/main.js":                     mainDotJS,
		"frontend/v1/public/index.html":               indexDotHTML,
		"frontend/v1/src/App.svelte":                  appDotSvelte,
		"frontend/v1/src/pages/Register.svelte":       registerDotSvelte,
		"frontend/v1/src/pages/Login.svelte":          loginDotSvelte,
		"frontend/v1/src/pages/ChangePassword.svelte": changePasswordDotSvelte,
		"frontend/v1/src/pages/Home.svelte":           homeDotSvelte,
		"frontend/v1/src/components/Table.svelte":     tableDotSvelte,
	}

	for _, typ := range types {
		files[fmt.Sprintf("frontend/v1/src/pages/%s/List.svelte", typ.Name.PluralRouteName())] = listDotSvelte(typ)
		files[fmt.Sprintf("frontend/v1/src/pages/%s/Create.svelte", typ.Name.PluralRouteName())] = createDotSvelte(typ)
		files[fmt.Sprintf("frontend/v1/src/pages/%s/Read.svelte", typ.Name.PluralRouteName())] = readDotSvelte(typ)
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(pkgRoot, filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return err
		}

		bytes := file()
		if _, err := f.Write(bytes); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func mainDotJS() []byte {
	return []byte(`import App from './App.svelte';

const app = new App({
	target: document.body,
	props: {
		// name: 'app'
	},
});

export default app;
`)
}

func indexDotHTML() []byte {
	f := []byte(`<!doctype html>
<html>
<head>
	<meta charset='utf8'>
	<meta name='viewport' content='width=device-width'>

	<title>App</title>

	<link rel='icon' type='image/png' href='favicon.ico'>
	<link rel='stylesheet' href='/global.css'>
	<link rel='stylesheet' href='/bundle.css'>
</head>

<body>
	<script src='/bundle.js'></script>
</body>
</html>
`)

	return f
}

func appDotSvelte() []byte {
	return []byte(`<script>
  import { Router, Link, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";

  // Auth routes
  import Login from "./pages/Login.svelte";
  import Register from "./pages/Register.svelte";
  import ChangePassword from "./pages/ChangePassword.svelte";

  /* // Items routes                                       */
  /* import ReadItem from "./pages/items/Read.svelte";     */
  /* import CreateItem from "./pages/items/Create.svelte"; */
  /* import Items from "./pages/items/List.svelte";        */

  export let url = "";
</script>

<!-- App.svelte -->

<Router {url}>
  <nav style="text-align: center;">
    <Link to="/">Home Page</Link>
    <!-- <Link to="items">Items</Link>                       -->
    <!-- <Link to="items/new">Create Item</Link>             -->
    <Link to="webhooks">Webhooks</Link>
    <Link to="login">Login</Link>
    <Link to="register">Register</Link>
    <Link to="password/new">Change Password</Link>
  </nav>
  <div>
    <!-- <Route path="items" component={Items} />            -->
    <!-- <Route path="items/:id" component={ReadItem} />     -->
    <!-- <Route path="items/new" component={CreateItem} />   -->
    <Route path="login" component={Login} />
    <Route path="register" component={Register} />
    <Route path="password/new" component={ChangePassword} />
    <Route path="/">
      <Home />
    </Route>
  </div>
</Router>
`)
}

func registerDotSvelte() []byte {
	return []byte(`<script>
  import { Link, navigate } from "svelte-routing";

  let username = "";
  let password = "";
  let passwordCopy = "";
  let twoFactorQRCode = "";

  // state vars
  let showingSecret = false;
  let canSubmit = false;

  function evaluateSubmission() {
    canSubmit =
      password.length > 0 && username.length > 0 && passwordCopy == password;
  }

  function moseyOn() {
    navigate("/login", { replace: true });
  }

  function handleRegistration() {
    fetch("/users/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        username,
        password
      })
    }).then(response => {
        if (response.status == 201) {
          return response.json();
        } else {
          console.error("something has gone awry: ");
          console.log(response);
        }
      }).then(data => {
        twoFactorQRCode = data["qr_code"];
        showingSecret = true;
      });
  }
</script>

<div style="margin-top: 7.5%; text-align: center;">

  {#if !showingSecret}
    <form id="registrationForm" on:submit|preventDefault={handleRegistration}>

      <p>
        username:
        <input
          bind:value={username}
          on:keyup={evaluateSubmission}
          type="text"
          name="username" />
      </p>

      <p>
        password:
        <input
          bind:value={password}
          on:keyup={evaluateSubmission}
          type="password"
          name="password" />
      </p>

      <p>
        once more so you're certain:
        <input
          bind:value={passwordCopy}
          on:keyup={evaluateSubmission}
          type="password" />
      </p>

      <input type="submit" value="register" disabled={!canSubmit} />
      <Link to="/login">log in instead</Link>

    </form>
  {:else}
    <img
      style="width: 20%;"
      src={twoFactorQRCode}
      alt="two factor authentication secret encoded as a QR code" />
    <p>
      You should save the secret this QR code contains, you'll be required to
      generate a token from it on every login.
    </p>
    <button on:click={moseyOn}>I've saved it, I promise</button>
  {/if}

</div>
`)
}

func loginDotSvelte() []byte {
	return []byte(`<script>
  import { Link, navigate } from "svelte-routing";

  let username = "";
  let password = "";
  let totp_token = "";
  let canSubmit = false;

  function evaluateSubmission() {
    canSubmit =
      password.length > 0 && username.length > 0 && totp_token.length > 0;
  }

  function handleLogin() {
    fetch("/users/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        username,
        password,
        totp_token
      })
    }).then(function(response) {
      if (response.status != 204) {
        console.error("something has gone awry");
      } else {
        console.log("login request was good");

        window.location.replace("/");
      }
    });
  }
</script>

<form
  id="loginForm"
  on:submit|preventDefault={handleLogin}
  style="margin-top: 7.5%; text-align: center;">
  <p>
    username:
    <input
      bind:value={username}
      on:keyup={evaluateSubmission}
      type="text"
      name="username" />
  </p>
  <p>
    password:
    <input
      bind:value={password}
      on:keyup={evaluateSubmission}
      type="password"
      name="password" />
  </p>
  <p>
    2FA code:
    <input bind:value={totp_token} on:keyup={evaluateSubmission} type="text" />
  </p>
  <input id="loginButton" type="submit" value="login" disabled={!canSubmit} />
  <Link to="/register">register instead</Link>
</form>
`)
}

func changePasswordDotSvelte() []byte {
	return []byte(`<script>
  import { Link, navigate } from "svelte-routing";

  let currentPassword = "";
  let newPassword = "";
  let newPasswordRepeat = "";
  let totpToken = "";
  let canSubmit = false;

  function evaluateSubmission() {
    canSubmit =
      newPassword.length > 0 &&
      currentPassword.length > 0 &&
      totpToken.length > 0 &&
      newPassword !== currentPassword &&
      newPasswordRepeat === newPassword;
  }

  function submitChangeRequest() {
    fetch("/users/password/new", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        new_password: newPassword,
        current_password: currentPassword,
        totp_token: totpToken
      })
    }).then(function(response) {
      if (response.status != 202) {
        console.error("something has gone awry");
      } else {
        navigate("/login", { replace: true });
      }
    });
  }
</script>

<form
  id="loginForm"
  on:submit|preventDefault={submitChangeRequest}
  style="margin-top: 7.5%; text-align: center;">
  <p>
    current password:
    <input
      bind:value={currentPassword}
      on:keyup={evaluateSubmission}
      type="password"
      name="username" />
  </p>
  <p>
    new password:
    <input
      bind:value={newPassword}
      on:keyup={evaluateSubmission}
      type="password"
      name="password" />
  </p>
  <p>
    once more so you're sure:
    <input
      bind:value={newPasswordRepeat}
      on:keyup={evaluateSubmission}
      type="password" />
  </p>
  <p>
    2FA code:
    <input bind:value={totpToken} on:keyup={evaluateSubmission} type="text" />
  </p>
  <input type="submit" value="change password" disabled={!canSubmit} />
</form>
`)
}

func homeDotSvelte() []byte {
	return []byte(`🆗`)
}

func tableDotSvelte() []byte {
	return []byte(`<script>
  export let columns = [];
  export let rows = [];
  export let tableStyle = "";
  export let rowClickFunc = () => {};
  export let rowDeleteFunc = () => {};

  let sortOrder = 1;
  let sortKey = "";
  let sortBy = r => "";
  let filterSettings = {};
  let columnByKey = {};

  columns.forEach(col => {
    columnByKey[col.key] = col;
  });

  $: c_rows = rows
    .filter(r =>
      Object.keys(filterSettings).every(f => {
        return (
          filterSettings[f] === undefined ||
          filterSettings[f] === columnByKey[f].filterValue(r)
        );
      })
    )
    .map(r => {
      return { ...r, $sortOn: sortBy(r) };
    })
    .sort((a, b) => {
      if (a.$sortOn > b.$sortOn) return sortOrder;
      else if (a.$sortOn < b.$sortOn) return -sortOrder;
      return 0;
    });

  const handleSort = col => {
    if (!col.unsortable) {
      if (sortKey === col.key) {
        sortOrder = sortOrder === 1 ? -1 : 1;
      } else {
        sortOrder = 1;
        sortKey = col.key;
        sortBy = r => r[sortKey];
      }
    }
  };
</script>

<style>
  .isSortable {
    cursor: pointer;
  }
</style>

<!-- heavily borrowed from/inspired by https://github.com/dasDaniel/svelte-table/blob/402a9eb3803ae2367f19651bddcb26ff46d29601/src/SvelteTable.svelte -->

{#if columns.length === 0}
  <h4>no data available :(</h4>
{:else}
  <table style={tableStyle}>
    <tr>
      {#each columns as col}
        <th
          on:click={() => handleSort(col)}
          class={col.sortable ? 'isSortable' : ''}>
        {col.title}
        {#if sortKey === col.key}{sortOrder === 1 ? '▲' : '▼'}{/if}
        </th>
      {/each}
      <th>
        <!--🗑 -->
      </th>
    </tr>
    {#each c_rows as row}
      <tr on:click={() => rowClickFunc(row)} style={row._style || ''}>
        {#each columns as col}
          <td>
            {@html col.renderValue ? col.renderValue(row) : row[col.key]}
          </td>
        {/each}
        <td>
          <button on:click={() => rowDeleteFunc(row)}>🗑️</button>
        </td>
      </tr>
    {/each}
  </table>
{/if}`)
}
