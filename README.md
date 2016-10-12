# PandaGo 🐼

[![Panda! Go, Panda!](http://i.imgur.com/rL88eG5.jpg)](https://youtu.be/nwd9W6VwIu4)

Pandoc, wrapped in Go.

## Cloning this Project

It's better to use `go get` instead of `git clone` so that your repo ends up
somewhere in the `$GOPATH` (which is relied on by `govendor`).

``` bash
go get github.com/scripted/pandago
cd $GOPATH/src/github.com/scripted/pandago
```

## PandaGo Locally

You can run the server directly.

``` bash
go run main.go
```

Compile and Run a la Heroku

``` bash
go install -v ./...
heroku local
```

Or run the app directly using `pandago`

``` bash
go install -v ./...
pandago
```

Ping root to see if it's working

``` bash
curl -i localhost:8080
#
# HTTP/1.1 200 OK
# Content-Type: application/json; charset=utf-8
# Date: Tue, 11 Oct 2016 21:53:14 GMT
# Content-Length: 19
#
# {"message":"OK 🐼, Go!"}
```

## Dependencies with govendor

PandaGo uses [`govendor`](https://github.com/kardianos/govendor) to manage its dependencies.

``` bash
go get -u github.com/kardianos/govendor
```

You can install all dependencies described in `vendor/vendor.json` into the
`vendor/` directory by running the following.

``` bash
govendor sync
```

To update dependencies as described in the app, run the following.

``` bash
govendor add +external
govendor remove +unused
```

For more info, checkout out [Heroku's page on govendor](https://devcenter.heroku.com/articles/go-dependencies-via-govendor).

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
curl -i https://<YOUR_APP_NAME>.herokuapp.com
#
# HTTP/1.1 200 OK
# Server: Cowboy
# Connection: keep-alive
# Content-Type: application/json; charset=utf-8
# Date: Tue, 11 Oct 2016 21:51:06 GMT
# Content-Length: 19
# Via: 1.1 vegur
#
# {"message":"OK 🐼, Go!"}
```

## Compiling Pandoc

There might be a time we need to recompile Pandoc. If that day comes, here are
some steps and tips.

### Heroku Cedar-14 Environment

You can get a fresh Cedar-14 environment using Docker.

``` bash
docker -it -v ~/local-shared-dir/:/container-shared-dir/ heroku/cedar:14 bash
```

You should probably go ahead and `apt-get update`

### Pandoc Dependencies

Pandoc needs [`hsb2hs`](https://hackage.haskell.org/package/hsb2hs) so we'll
install [Haskell Platform](https://www.haskell.org/platform/) to get it using
`cabal`. Then, append the Cabal binary path to `$PATH` so `hsb2hs` can be found.

``` bash
apt-get update
apt-get install haskell-platform
cabal update
cabal install hsb2hs
export PATH=$PATH:/root/.cabal/bin/
```

We'll also install [Haskell Stack](https://docs.haskellstack.org/en/stable/README/),
which we will use to actually build Pandoc.

```
curl -sSL https://get.haskellstack.org/ | sh
stack setup
```

### Compiling Pandoc

Check the `INSTALL` instructions included in the Pandoc source.

We'll use Haskell Stack to install Pandoc, making sure to embed data-files so
that the binary is relocatable. These include template files among others that
Pandoc needs to create `.docx` files, for example.

``` bash
wget https://hackage.haskell.org/package/pandoc-1.17.0.3/pandoc-1.17.0.3.tar.gz
tar xvzf pandoc-1.17.0.3.tar.gz
cd pandoc-1.17.0.3
stack install pandoc --flag pandoc:embed_data_files
```

### Tarballing and Deploying

The [Custom Binaries buildpack](https://github.com/tonyta/heroku-buildpack-custom-binaries)
expects a tarball.

``` bash
cp /root/.local/bin/pandoc /container-shared-dir/
cd /container-shared-dir/
tar -cvzf pandoc-embedded.tar.gz pandoc
```

Then we can just find `pandoc-embedded.tar.gz` at `~/local-shared-dir/` on our
local machine.

From here, all we have to do is host it publicly (e.g, on AWS S3) and update our
`.custom_binaries` manifest.
