# golang:1.20.2-alpine3.17 SHA1 digest.
FROM golang@sha256:110b07af87238fbdc5f1df52b00927cf58ce3de358eeeb1854f10a8b5e5e1411 as builder

ARG BIN

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

ENV USER=${BIN}
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nohome" \
    --no-create-home \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"


FROM scratch

ARG BIN
ARG PORT

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY ./build/${BIN} /${BIN}

EXPOSE ${PORT}

USER ${BIN}:${BIN}

ENTRYPOINT [ "/endpoint" ]
