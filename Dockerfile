FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY gppMonitor /app/
ENTRYPOINT ["/app/gppMonitor"]
EXPOSE 8001