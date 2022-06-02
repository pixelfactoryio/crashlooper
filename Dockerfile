FROM golang:1.18 as builder
WORKDIR /build
COPY . .
RUN go mod download
RUN make

FROM debian:latest
EXPOSE 3000
COPY --from=builder /build/bin/crashlooper /usr/local/bin/crashlooper
ENTRYPOINT [ "/usr/local/bin/crashlooper" ]