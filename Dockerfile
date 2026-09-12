FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM gcr.io/distroless/static:nonroot
COPY --from=build /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
