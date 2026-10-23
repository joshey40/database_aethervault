# syntax=docker/dockerfile:1

# build stage
FROM golang:alpine AS build
WORKDIR /src

# copy only the go mod and go sum and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy the whole app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/aethervault

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS deploy

WORKDIR /app/

COPY --from=build /bin/aethervault aethervault
COPY --from=build /src/config.yml config.yml

EXPOSE 3333

USER nonroot:nonroot

ENTRYPOINT ["/app/aethervault"]