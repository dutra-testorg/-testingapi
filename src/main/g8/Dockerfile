# STEP 1 build executable binary
FROM golang:1.19 as builder
LABEL stage=intermediate

COPY . /$name;format="lower,hyphen"$/
WORKDIR /$name;format="lower,hyphen"$/

ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN ssh-keyscan github.com > /root/.ssh/known_hosts
RUN echo "\${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN go env -w GOPRIVATE="github.com/Gympass/*,github.com/gympass/*"
RUN git fetch --all --tags

RUN make deps
RUN make build

# STEP 2 build a small image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /$name;format="lower,hyphen"$/bin/app /$name;format="lower,hyphen"$

EXPOSE 8080

ENTRYPOINT [ "/$name;format="lower,hyphen"$" ]
