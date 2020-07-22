echo "COMPILING!"
cd server_api
go build .
cd ../worker_test
go build .
cd ../command_logger
go build .
echo "MOVING!"
cd ..
rm -rf buildmkdir build
mkdir build
mv server_api/server_api ./build/
mv worker_test/worker_test ./build/
mv command_logger/command_logger ./build/
echo "DONE!"
echo "sudo -b -E ./command_logger ./server_api" >> ./build/run_server
chmod +x ./build/run_server
echo "sudo -b -E ./command_logger ./worker_test" >> ./build/run_worker
chmod +x ./build/run_worker
cp db/sql/schema/schema.sql ./build/
