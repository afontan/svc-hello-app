# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o /out/svc-hello .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/svc-hello /svc-hello
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/svc-hello"]
