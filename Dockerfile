# syntax=docker/dockerfile:1.2

ARG GO_VERSION=1.17
ARG GOLANGCI_LINT_VERSION=v1.43.0

FROM golang:${GO_VERSION}-stretch AS base
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

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    make -f rules/compiler-builder.mk check

FROM base as make-compiler

ARG DEBUG

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    APP_EXECUTABLE_OUT=/out \
    DEBUG=${DEBUG} \
    make -f rules/compiler-builder.mk build

FROM scratch as compiler-out
COPY --from=make-compiler /out/* .

FROM base AS make-go-mod-tidy
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy

FROM scratch AS go-mod-tidy
COPY --from=make-go-mod-tidy /app/go.mod .
COPY --from=make-go-mod-tidy /app/go.sum .

FROM base as make-lexer

RUN --mount=target=. \
    APP_EXECUTABLE_OUT=/out \
    make -f rules/lexer-builder.mk build

FROM scratch as lexer-out
COPY --from=make-lexer /out/* .

FROM debian:9 as compiler-image

ARG LEXER_PATH=/app/bin/lexer

ENV LEXER_PATH=${LEXER_PATH}

COPY --from=lexer-out lexer /app/bin/
COPY --from=compiler-out compiler /app/bin/

ENTRYPOINT /app/bin/compiler
