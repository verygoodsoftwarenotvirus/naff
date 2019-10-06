#! /bin/bash

for dir in `go list gitlab.com/verygoodsoftwarenotvirus/todo/...`; do rm workspaces/$(echo "${dir/gitlab.com\/verygoodsoftwarenotvirus\/todo\/}" | sed -r 's/\//_/g').code-workspace; done
