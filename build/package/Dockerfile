## Build stage
FROM golang:1.12.5 as build

WORKDIR /src/k8s-device-plugin-socketcan

COPY *.mod ./
RUN go mod download

COPY . .

# Build and install
RUN make install

## Run stage
FROM golang:1.12.5
COPY --from=build /go/bin/socketcan /bin/socketcan

ENTRYPOINT ["/bin/socketcan"]