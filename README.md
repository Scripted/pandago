# Go-Panda

Pandoc, wrapped in Go

## Cloning this Project

It's better to use `go get` instead of `git clone` so that your repo ends up
somewhere in the `$GOPATH` (which is relied on by `govendor`).

``` bash
go get github.com/scripted/go-panda
```

## Run Locally

``` bash
go run main.go
```

## Deploying to Heroku

Create the Heroku app as usual:

``` bash
heroku create
```

Add the official `heroku/go` buildpack as well as the a buildpack for installing
custom binaries so we can install Pandoc via `.custom_binaries`:

``` bash
heroku buildpacks:add heroku/go
heroku buildpacks:add https://github.com/tonyta/heroku-buildpack-custom-binaries#v1.0.0
```

That's it. Then just push it to Heroku like normal.

``` bash
git push heroku master
curl -i https://go-panda.herokuapp.com/HTTP/1.1 200 OK
# Server: Cowboy
# Connection: keep-alive
# Content-Type: application/json; charset=utf-8
# Date: Tue, 11 Oct 2016 21:20:56 GMT
# Content-Length: 19
# Via: 1.1 vegur
#
# {"message":"pong"}
```
