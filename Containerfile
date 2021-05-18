FROM registry.access.redhat.com/ubi8/ubi as build
LABEL maintainer=teamnado-ops@redhat.com

ENV GODIR=/usr/local/go APPDIR=/opt/app CGO_ENABLED=0

RUN curl -sfL --retry 10 -o /tmp/go.tar.gz https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && \
    echo "7154e88f5a8047aad4b80ebace58a059e36e7e2e4eb3b383127a28c711b4ff59 /tmp/go.tar.gz" | sha256sum -c && \
    mkdir -p $GODIR && \
    tar --strip-components=1 -zxf /tmp/go.tar.gz --directory $GODIR && \
    rm /tmp/go.tar.gz && \
    mkdir -p $APPDIR

ENV PATH="$GODIR/bin:$PATH"

WORKDIR $APPDIR

COPY . $APPDIR

RUN go build -o caddy cmd/caddy/caddy.go

# Multi-stage
FROM registry.access.redhat.com/ubi8/ubi-minimal as run

COPY misc/rhit.crt /etc/pki/ca-trust/source/anchors/rhit.crt
RUN update-ca-trust

WORKDIR /opt/app/

COPY --from=build /opt/app/caddy .

CMD /opt/app/caddy
