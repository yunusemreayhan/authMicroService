FROM nginx:alpine

RUN mkdir -p /public

COPY ../../public /public

COPY ./proxy/nginx_test.conf /etc/nginx/conf.d/default.conf