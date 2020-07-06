http conf:
-----------
`basepath` in `front/src/routes/utils.js`
port, CORS and others in `back/server_api/serve.go`
redirect and homepage in github oauth app
POSTGRES conf in `back/db/dbconn.go`
env variables: `TD4_*`

Frontend
-----------
npm run dev
npm build
npx sapper export
npx serve __sapper__/export

SQLC
------------
sqlc generate

useful stuff:
----------------
check which process is using a port: `lsof -i :8081`
go linter being used: golangci-lint: https://github.com/golangci/golangci-lint#binary
screen: `gtk-redshift -l 32.08:34.78`


tech used:
--------------
pytest
postgresql
golang
sapper (...svelte)
docker
sqlc
api + static rendered site
github pages
vultr


update packages
---------------------
go:

in the package root folder:
go get -u 

upgrade sqlc:
go get -u github.com/kyleconroy/sqlc/cmd/sqlc

npm:

edit package.json for new versions
npm update --save
npm update --save-dev


# TODO: mytests
# TODO: mysolutions
# TODO: content in "about" page
# TODO: deletion of tests
# TODO: deletion of solutions
# TODO: global stats
# TODO: nicer ui
# TODO: separate static path (in fronetend!)
# TODO: invalidate cache on new test and make cache TTL longer
# TODO: a "loading..." message instead of "no tests yet". These are different cases!
# TODO: test page- when reloaded is all undefined!
# TODO: remove old, stopped containers periodically
# TODO: cache only valid urls



