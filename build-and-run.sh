#!/bin/bash

cd client || exit

npm ci
npm run build

cd ../server

go build -tags netgo -ldflags '-s -w' -o ./bin/app ./cmd/api/.

./bin/app
