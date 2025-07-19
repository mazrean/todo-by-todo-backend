# syntax=docker/dockerfile:1
FROM golang:1.24.1 AS build

WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /bin/app /bin/app

ENTRYPOINT ["/bin/app"]
