http conf:
-----------
`basepath` in `front/src/utils.js`
port and CORS in `back/serve.go`
redirect and homepage in github oauth app
POSTGRES conf in `sql/db/dbconn.go`
env variable: `TD4`

Frontend
-----------
npm run dev
npm build
npx sapper export
npx serve __sapper__/export

SQLC
------------
sqlc generate

server
----------
./dev_run.py

worker
-----------
currently run manually: `go run .` or `go run worker` etc.

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
svelte
sapper
docker
sqlc
api + static rendered site


update packages
---------------------
go get -u all
edit package.json for new versions
npm update --save
npm update --save-dev




# TODO: implement new_solution page
# TODO: limit characters in test / solution desc
# TODO: limit code in test
# TODO: limit code in solution
# TODO: limit title length
