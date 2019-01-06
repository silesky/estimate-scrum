# base image
FROM node:carbon-alpine

# set working directory
RUN mkdir /app
WORKDIR /app

# add `/usr/src/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install and cache app dependencies.
# If package.json or package-lock changes, `docker-compose build` will automatically trigger a fresh npm install
COPY package.json package-lock.json app/

RUN npm install --silent
RUN npm install react-scripts@1.1.1 -g --silent

# start app
CMD ["npm", "start"]
