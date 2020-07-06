npx sapper export
rm -rf build
mkdir build
cp -r __sapper__/export/* ./build/
echo "solvemytest.dev" > ./build/CNAME
cd ../../td4_front
rm * -rf
cp -r ../td4/front/build/* .
git add .
git commit -am "new build"
echo "DONE!"
