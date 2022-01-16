# syntax=docker/dockerfile:1.2

ARG GO_VERSION=1.17
ARG GOLANGCI_LINT_VERSION=v1.43.0
ARG DENO_VERSION=1.16.4

FROM golang:${GO_VERSION}-stretch AS go-base
WORKDIR /app

# Add repos with clang
RUN echo "deb http://archive.ubuntu.com/ubuntu bionic main multiverse restricted universe" >> /etc/apt/sources.list && \
    echo "deb http://archive.ubuntu.com/ubuntu bionic-security main multiverse restricted universe" >> /etc/apt/sources.list && \
    echo "deb http://archive.ubuntu.com/ubuntu bionic-updates main multiverse restricted universe" >> /etc/apt/sources.list

RUN apt-get update && apt-get install -y \
    # To be able to download clang
    --allow-unauthenticated \
    make \
    clang-6.0

COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} AS lint-base

FROM debian:9 as debian-base

FROM go-base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    make -f rules/compiler-builder.mk check

FROM go-base as make-compiler

ARG DEBUG

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    APP_EXECUTABLE_OUT=/out \
    DEBUG=${DEBUG} \
    make -f rules/compiler-builder.mk build

FROM scratch as compiler-out
COPY --from=make-compiler /out/* .

FROM go-base AS make-go-mod-tidy
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy

FROM scratch AS go-mod-tidy
COPY --from=make-go-mod-tidy /app/go.mod .
COPY --from=make-go-mod-tidy /app/go.sum .

FROM go-base as make-lexer

RUN --mount=target=. \
    APP_EXECUTABLE_OUT=/out \
    make -f rules/lexer-builder.mk build

FROM scratch as lexer-out
COPY --from=make-lexer /out/* .

FROM denoland/deno:${DENO_VERSION} as deno-base

WORKDIR /app

RUN apt-get update && apt-get install -y \
    make

FROM deno-base as make-astbackend

RUN --mount=target=. \
    APP_EXECUTABLE_OUT=/out \
    make -f rules/astbackend-builder.mk build

FROM scratch as astbackend-out
COPY --from=make-astbackend /out/* .

FROM debian-base as compiler-image

ARG LEXER_PATH=/app/bin/lexer
ARG ASTBACKEND_PATH=/app/bin/astbackend

ENV LEXER_PATH=${LEXER_PATH}
ENV ASTBACKEND_PATH=${ASTBACKEND_PATH}

COPY --from=lexer-out lexer /app/bin/
COPY --from=compiler-out compiler /app/bin/
COPY --from=astbackend-out astbackend /app/bin/

ENTRYPOINT /app/bin/compiler
