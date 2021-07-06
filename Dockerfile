FROM golang:1.16.5-alpine AS build

WORKDIR /go/src/kiddy-lp

COPY . .

RUN go install ./...

FROM alpine:3.12
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD cmd

#docker build . -t kiddy
#docker run --rm --link redis --link lp -p 8080:8080 -p 8081:8081 --env STORAGE="redis:6379" --env LP_ADDRESS="http://lp:8000"  -it --name kiddy kiddy cmd
#docker kill kiddy