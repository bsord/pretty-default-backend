FROM nginx:1.17.10

# NGINX Config
COPY ./nginx.conf /etc/nginx/nginx.conf

# Resources
COPY content/ /var/www/html/

CMD ["nginx", "-g", "daemon off;"]