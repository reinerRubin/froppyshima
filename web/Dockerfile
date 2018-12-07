FROM nginx:stable

COPY ./docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY . /var/www

CMD /usr/sbin/nginx -g "daemon off;"
