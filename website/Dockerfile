# https://cli.vuejs.org/zh/guide/deployment.html#docker-nginx
# build stage
FROM node:lts-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf

# docker run -d -p 8080:80 txfs19260817/scopelens-website
