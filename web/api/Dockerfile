# ---- Base Node ----
FROM golang:1.12-alpine3.10 AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
WORKDIR /app/api
RUN apk update && apk upgrade && apk add --no-cache git
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
COPY --from=build /app/api/apiapp ./
RUN chmod +x /app/apiapp
EXPOSE 80
CMD ./apiapp
