
ARG golang_version=1.23

FROM golang:${golang_version}-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags "-s -w" -v -o /app/main ./cmd/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/main /main

CMD ["/main"]
