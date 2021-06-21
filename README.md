### About

Project with compiler written on golang and lexer written on C++

### Dependencies

Languages:
* [Golang](https://golang.org/)
* [C++ - clang compiler](https://clang.llvm.org/)

Instruments:
* [Docker](https://www.docker.com/)
* [Make](https://ru.wikipedia.org/wiki/Make)


### Build and run

To run slr-runner run:
* `make` - it will build everything and run tests
* `make run-slr-runner` - it`s run slr-runner with some predefined params to easy rerun for development 

**If you don`t have clang with required version you can use docker to build binaries**

To build in docker:
* Run `make build-dproxy` - it`s builds builder container
* That`s all

**To build everything you now can run**
`./bin/dmake build`

**Also, when you use dmake available other build steps from Makefile**