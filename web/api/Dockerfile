# ---- Base Node ----
FROM golang:1.13beta1-alpine3.10 AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
WORKDIR /app/api
RUN apk update && apk upgrade && apk add --no-cache git ca-certificates && update-ca-certificates
COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download

# ---- Copy Files/Build ----
FROM dependencies AS build
WORKDIR /app
COPY ./api /app/api
WORKDIR /app/api
ENV ENV=production
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apiapp

# --- Release with Alpine ----
FROM alpine:3.10 AS release
WORKDIR /app
COPY --from=dependencies /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/api/migrations ./migrations
COPY --from=build /app/api/assets ./assets
COPY --from=build /app/api/apiapp ./
RUN chmod +x /app/apiapp
EXPOSE 80
CMD ./apiapp
