package frontendsrc

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func listDotSvelte(typ models.DataType) func() []byte {
	f := `<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },`

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	cn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()
	uvn := typ.Name.UnexportedVarName()

	for _, field := range typ.Fields {
		f += fmt.Sprintf(`
    {
      title: %q,
      key: %q
    },`, field.Name.Singular(), field.Name.RouteName())
	}
	f += fmt.Sprintf(`
    {
      title: "Created On",
      key: "created_on"
    },
    {
      title: "Updated On",
      key: "updated_on"
    }
  ];
  let %s = [];

  function delete%s(row) {
    if (confirm("are you sure you want to delete this %s?")) {
      fetch(`+"`"+`/api/v1/%s/${row.id}`+"`"+`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        %s = %s.filter(%s => {
          return %s.id != row.id;
        });
      });
    }
  }
`, puvn, sn, cn, prn, puvn, puvn, uvn, uvn)

	f += fmt.Sprintf(`
  function goTo%s(row) {
    navigate(`+"`"+`/%s/${row.id}`+"`"+`, { replace: true });
  }

  fetch("/api/v1/%s/")
    .then(response => response.json())
    .then(data => {
      %s = data["%s"];
    });
</script>

<!-- %s.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ %s }
  rowClickFunc={goTo%s}
  rowDeleteFunc={delete%s} />`, sn, prn, prn, puvn, prn, pn, puvn, sn, sn)

	return func() []byte { return []byte(f) }
}

func createDotSvelte(typ models.DataType) func() []byte {

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	prn := typ.Name.PluralRouteName()
	uvn := typ.Name.UnexportedVarName()

	f := fmt.Sprintf(`<script>
  import { Link } from "svelte-routing";

  let name = "";
  let details = "";
  let canSubmit = false;

  function evaluateSubmission() {
    canSubmit = name != "" && details != "";
  }

  function create%s() {
    fetch("http://localhost/api/v1/%s/", {
      method: "POST",
      mode: "cors", // no-cors, cors, *same-origin
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        name,
        details
      })
    }).then(function(response) {
      if (response.status != 201) {
        console.error("something has gone awry");
      } else {
        name = "";
        details = "";
      }
    });
  }
</script>

<!-- %s.svelte -->
<form
  id="%sForm"
  on:submit|preventDefault={create%s}
  style="margin-top: 7.5%%; text-align: center;">
  <p>
    name:
    <input
      bind:value={name}
      on:keyup={evaluateSubmission}
      type="text"
      name="name" />
  </p>
  <p>
    details:
    <input
      bind:value={details}
      on:keyup={evaluateSubmission}
      type="text"
      name="details" />
  </p>
  <input type="submit" value="create" disabled={!canSubmit} />
  <Link to="/%s">%s list</Link>
</form>
`, sn, prn, pn, uvn, sn, prn, pcn)

	return func() []byte { return []byte(f) }
}

func readDotSvelte(typ models.DataType) func() []byte {
	pn := typ.Name.Plural()
	rn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()
	uvn := typ.Name.UnexportedVarName()

	f := fmt.Sprintf(`<script>
  import { Link } from "svelte-routing";

  let %s = {};

  const %sID = window.location.pathname.replace("/%s/", "");

  fetch(`+"`"+`/api/v1/%s/${ %sID }`+"`"+`)
    .then(response => response.json())
    .then(data => {
      %s = data;
    });
</script>

<!-- %s.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html %s.name}
  </p>
  <p>
    details:
    {@html %s.details}
  </p>

  <Link to="/%s">see all</Link>

</div>`, uvn, uvn, prn, prn, uvn, uvn, pn, rn, rn, prn)

	return func() []byte { return []byte(f) }
}
