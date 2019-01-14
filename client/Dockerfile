# base image
FROM node:carbon-alpine

# set working directory
RUN mkdir /app
WORKDIR /app

# add `/usr/src/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# Copy package.json package-lock from local to working directory
COPY package.json package-lock.json ./

# If package.json or package-lock changes, `docker-compose build` will automatically trigger a fresh npm install
# install and cache app dependencies.
RUN npm install
RUN npm install react-scripts@1.1.1 -g --silent

# start app
CMD ["npm", "start"]
