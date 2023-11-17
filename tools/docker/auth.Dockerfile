FROM golang:bullseye as runner

# Add the binary
COPY ./build/auth_micro_service_server /bin/auth_micro_service_server
COPY ./config.json /config.json

COPY ./build/private.key /private.key
COPY ./internal/key/myCA.crt /bin/myCA.crt
COPY ./internal/key/myCA.key /bin/myCA.key
COPY ./internal/key/mycert1.crt /bin/mycert1.crt
COPY ./internal/key/mycert1.key /bin/mycert1.key
COPY ./internal/key/mycert1.req /bin/mycert1.req
COPY ./internal/key/myclient.crt /bin/myclient.crt
COPY ./internal/key/myclient.key /bin/myclient.key
COPY ./internal/key/myclient.csr /bin/myclient.csr
COPY ./internal/key/myclient.p12 /bin/myclient.p12

# Run the auth_micro_service command by default when the container starts.
CMD cd /bin && ./auth_micro_service_server --tls-certificate mycert1.crt --tls-key mycert1.key --host 0.0.0.0 --port 3000
