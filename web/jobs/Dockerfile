# ---- Base Node ----
FROM golang:1.13.1-alpine3.10 AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
WORKDIR /app/jobs
RUN apk update && apk upgrade && apk add --no-cache git ca-certificates gcc cmake cmake-doc extra-cmake-modules extra-cmake-modules-doc build-base bash && update-ca-certificates
COPY ./jobs/go.mod ./
COPY ./jobs/go.sum ./
RUN go mod download

# ---- Copy Files/Build ----
FROM dependencies AS build
WORKDIR /app
COPY ./jobs /app/jobs
RUN mkdir /go/pkg/mod/github.com/google/brotli@v1.0.7/out
WORKDIR /go/pkg/mod/github.com/google/brotli@v1.0.7/out
RUN chmod +x ../configure-cmake
RUN ../configure-cmake --disable-debug
RUN make && make test && make test install
WORKDIR /app/jobs
ENV ENV=production
RUN sed -i '/#cgo/d' /go/pkg/mod/github.com/google/brotli@v1.0.7/go/cbrotli/cgo.go
RUN CGO_ENABLED=1 CGO_LDFLAGS="-lbrotlidec-static -lbrotlienc-static -lbrotlicommon-static -lm -static" GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -o jobsapp

# --- Release with Alpine ----
FROM alpine:3.10 AS release
WORKDIR /app
COPY --from=dependencies /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/jobs/playlists.json ./
COPY --from=build /app/jobs/jobsapp ./
RUN chmod +x /app/jobsapp
EXPOSE 80
CMD ./jobsapp
