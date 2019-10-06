#! /bin/bash

for dir in `go list gitlab.com/verygoodsoftwarenotvirus/todo/...`; do echo "{\"folders\":[{\"path\": \"/home/jeffrey/src/${dir/todo/naff\/example_output}\"},{\"path\": \"/home/jeffrey/src/$dir\"}]}" | jq . -M > workspaces/$(echo "${dir/gitlab.com\/verygoodsoftwarenotvirus\/todo\/}" | sed -r 's/\//_/g').code-workspace; done
