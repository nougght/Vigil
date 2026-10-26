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
COPY backend/services/core ./backend/services/core

WORKDIR /build/backend
RUN go build -o /build/monit-core ./services/core/cmd/main.go


# light weight image for running
FROM alpine
WORKDIR /app/core

COPY --from=build /build/monit-core .
COPY backend/services/core/config.yaml .
COPY bin/agent.exe /app/bin/agent.exe
COPY backend/services/core/creds/root-ca.crt /app/server/creds/root-ca.crt
ENTRYPOINT ["./monit-core"]
