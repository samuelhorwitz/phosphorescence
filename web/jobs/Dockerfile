# ---- Base Node ----
FROM golang:1.13beta1-alpine3.10 AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
WORKDIR /app/jobs
RUN apk update && apk upgrade && apk add --no-cache git ca-certificates && update-ca-certificates
COPY ./jobs/go.mod ./
COPY ./jobs/go.sum ./
RUN go mod download

# ---- Copy Files/Build ----
FROM dependencies AS build
WORKDIR /app
COPY ./jobs /app/jobs
WORKDIR /app/jobs
ENV ENV=production
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jobsapp

# --- Release with Alpine ----
FROM alpine:3.10 AS release
WORKDIR /app
COPY --from=dependencies /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/jobs/playlists.json ./
COPY --from=build /app/jobs/jobsapp ./
RUN chmod +x /app/jobsapp
EXPOSE 80
CMD ./jobsapp
