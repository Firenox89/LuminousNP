#! /bin/bash

rm -rf build
rm -rf web/dist

mkdir build

pushd web
yarn build
popd
pushd backend
./crosscompile.sh
popd
cp -R web/dist build/
cp backend/nodemcu-controller build/

ssh pi@nodemcu-controller 'rm -rf /home/pi/dist'
scp build/nodemcu-controller pi@nodemcu-controller:/home/pi
scp -r build/dist pi@nodemcu-controller:/home/pi
