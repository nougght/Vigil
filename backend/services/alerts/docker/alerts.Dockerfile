FROM golang:1.26-alpine AS init

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
# copy only alerts service
COPY backend/services/alerts ./backend/services/alerts

WORKDIR /build/backend


FROM init AS build

RUN CGO_ENABLED=0 go build -o /build/monit-alerts ./services/alerts/cmd/main.go


FROM init AS build-debug

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /build/monit-alerts ./services/alerts/cmd/main.go


# light weight image for running
FROM alpine AS run
WORKDIR /app/alerts

COPY --from=build /build/monit-alerts .
ENTRYPOINT ["./monit-alerts"]


FROM alpine AS run-debug
WORKDIR /app/alerts

COPY --from=build-debug /build/monit-alerts .
COPY --from=build-debug /go/bin/dlv .
ENTRYPOINT ["./dlv", "exec", "./monit-alerts", "--headless", "--listen=:40012", "--api-version=2", "--accept-multiclient", "--continue"]
