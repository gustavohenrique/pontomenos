image="gustavohenrique/pontomenos"
#image="golang:1.8.3-alpine3.6"

docker run --rm -v ${PWD}:/code ${image} sh -c 'cd /code \
 && apk add --update git \
 && go get github.com/gin-gonic/gin \
 && go get github.com/gin-contrib/cors \
 && go get github.com/parnurzeal/gorequest \
 && go build main.go'

docker build -t=registry.heroku.com/pontomenos/web .
docker push registry.heroku.com/pontomenos/web

