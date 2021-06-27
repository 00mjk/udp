## Copyright (c) 2021 Abhijit Bose. All Right reserved.
## Use of this source code is governed by a Apache 2.0 license that can be found
## in the LICENSE file.

.PHONY: test, build, format, clean


format:
	go vet ./...
	go fmt ./...

test:
	go test -race -v .
	go clean -testcache


