echo "COMPILING!"
cd back/server_api
go build .
cd ../worker_test
go build .
echo "MOVING!"
cd ..
rm -rf buildmkdir build
mkdir build
mv server_api/server_api ./build/
mv worker_test/worker_test ./build/
echo "DONE!"