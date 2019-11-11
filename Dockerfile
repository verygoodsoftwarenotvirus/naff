FROM golang:latest AS build

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/naff

COPY . .

RUN go build -o /naff -ldflags "-X main.Version=$(git rev-parse --short HEAD)" gitlab.com/verygoodsoftwarenotvirus/naff/cmd/cli

FROM debian:latest

COPY --from=build /naff /naff

RUN ls -Al /naff

ENTRYPOINT /naff
