FROM docker.io/golang:1.21 AS build

WORKDIR /zssnd

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./zssnd ./cmd/zssnd

FROM busybox

COPY --from=build /zssnd/zssnd /zssnd

ENTRYPOINT ["/zssnd"]
