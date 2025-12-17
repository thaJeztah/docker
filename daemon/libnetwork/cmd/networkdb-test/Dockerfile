FROM dhi.io/alpine-base:3.22

RUN apk --no-cache add curl

COPY testMain /app/

WORKDIR app

ENTRYPOINT ["/app/testMain"]
