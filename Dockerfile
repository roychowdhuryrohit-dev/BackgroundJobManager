FROM golang:1.13-alpine AS builder-container
WORKDIR /app
ADD . /app
RUN apk add --update --no-cache git
RUN apk --update --no-cache add ca-certificates && \
    update-ca-certificates
# gcc musl-dev
RUN cd /app && mkdir -p /app/bin && \
    CGO_ENABLED=0 go build -o /app/bin/DemoService -tags netgo
# FROM alpine
# RUN apk update && \
#     apk add ca-certificates && \
#     update-ca-certificates && \
#     rm -rf /var/cache/apk/*
# WORKDIR /app
# COPY --from=builder-container /app/ProfileService /app
# EXPOSE 8080
# ENTRYPOINT ["./ProfileService"]
FROM scratch
WORKDIR /app

COPY --from=builder-container /app/bin/DemoService /app/bin/DemoService
COPY --from=builder-container /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /dump/baseline
ENV TEAM_CSV /dump/baseline

WORKDIR /dump/team
ENV BASELINE_CSV /dump/team

ENTRYPOINT [ "/app/bin/DemoService" ]
