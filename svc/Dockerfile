FROM frolvlad/alpine-glibc

RUN apk --no-cache add ca-certificates && update-ca-certificates

RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe


# Needed for skaffold to detect this is a debug container
ENV GOTRACEBACK=all

ARG binary

COPY ${binary} main

STOPSIGNAL SIGTERM

ENTRYPOINT ["./main"]


