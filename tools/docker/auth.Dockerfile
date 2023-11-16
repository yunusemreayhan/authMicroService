FROM golang:alpine as runner
FROM alpine:3.5

# Add the binary
COPY ./public/ /public
COPY ./build/auth_micro_service /bin/auth_micro_service
COPY ./config.json /config.json
COPY ./build/private.key /private.key

# Run the auth_micro_service command by default when the container starts.
CMD ["/bin/auth_micro_service"]
