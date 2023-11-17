FROM nginx:alpine

RUN mkdir -p /public

COPY ../../public /public

COPY ./proxy/nginx.conf /etc/nginx/conf.d/default.conf