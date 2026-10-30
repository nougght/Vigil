FROM golang:1.26-alpine AS build

WORKDIR /build

COPY backend/go.mod backend/go.sum ./backend/
# shared modules
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY gen/go.mod gen/go.sum ./gen/

RUN cd backend && go mod download -x

COPY pkg ./pkg
COPY gen ./gen
# shared services pkg
COPY backend/pkg ./backend/pkg
# copy only core service
COPY backend/services/alerts ./backend/services/alerts

WORKDIR /build/backend
RUN go build -o /build/monit-alerts ./services/alerts/cmd/main.go


# light weight image for running
FROM alpine
WORKDIR /app/alerts

COPY --from=build /build/monit-alerts .
ENTRYPOINT ["./monit-alerts"]
