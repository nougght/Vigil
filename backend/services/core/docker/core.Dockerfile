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
# copy only core service
COPY backend/services/core ./backend/services/core

WORKDIR /build/backend

# delve for debug mode
FROM golang:1.26-alpine AS dlv
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@v1.27.0

FROM init AS build
RUN CGO_ENABLED=0 go build -o /build/monit-core ./services/core/cmd/main.go

FROM init AS build-debug
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /build/monit-core ./services/core/cmd/main.go



# light weight image for running
FROM alpine AS run
WORKDIR /app/core

COPY --from=build /build/monit-core .
COPY backend/services/core/config.yaml .
COPY bin/agent.exe /app/bin/agent.exe
COPY backend/services/core/creds/root-ca.crt /app/server/creds/root-ca.crt
ENTRYPOINT ["./monit-core"]


FROM alpine AS run-debug
WORKDIR /app/core

COPY --from=build-debug /build/monit-core .
COPY --from=dlv /go/bin/dlv .
COPY backend/services/core/config.yaml .
COPY bin/agent.exe /app/bin/agent.exe
COPY backend/services/core/creds/root-ca.crt /app/server/creds/root-ca.crt
ENTRYPOINT ["./dlv", "exec", "./monit-core", "--headless", "--listen=:40011", "--api-version=2", "--accept-multiclient", "--continue"]
