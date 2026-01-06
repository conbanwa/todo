FROM golang:1.21-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/todo .

FROM alpine:3.18
COPY --from=build /bin/todo /bin/todo
EXPOSE 8080
ENTRYPOINT ["/bin/todo"]
