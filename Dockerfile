FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY GPPMonitor /app/
ENTRYPOINT ["/app/GPPMonitor"]
EXPOSE 8002