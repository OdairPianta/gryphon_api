# Pull nginx base image
FROM nginx:latest

# Expost port 80
EXPOSE 80

# Copy custom configuration file from the current directory
COPY ./docker/nginx.conf /etc/nginx/nginx.conf