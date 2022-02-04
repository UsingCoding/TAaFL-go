### Compiler - compiler from our specific language to JS

**Repository archived since work on project completed**

### About

#### What is

Compiler written in Go, uses Lexer written in c++ and AST backend with library [escodegen](https://github.com/estools/escodegen) - JS code generator

#### Why JS

Since JavaScript is popular language, and it works completely everywhere. So thus when you write program in our language it`s cross-platform out of the box.

### Dependencies

* [Golang](https://golang.org/)
* [C++ - clang compiler](https://clang.llvm.org/)
* [Deno - modern JS environment](https://clang.llvm.org/)
  * [escodegen](https://github.com/estools/escodegen) - AST code generator - library just copied into `/data/astbackend/escodegen` and adopted to use in Deno environment.  

### Usage

Compiler supports compiling only one file.

Since compiler distributing via Docker image you can run it by passing all argument to container startup command

```shell
docker run --rm -name compiler-script -i vadimmakerov/compiler < in.script > app.js
```

Input file passed via stdin and output printed to stdout so it can be redirected further for example JS optimizer or to file 

### Build

Tools to build compiler

- `make` - To aggregate docker build commands
- [Docker BuildKit](https://github.com/moby/buildkit) - Improved build system
- [Docker BuildX](https://github.com/docker/buildx) - docker build plugin

After you install all of it just run
```shell
make
```

It will run linter check and build Docker image with tag `vadimmakerov/compiler:master` 
