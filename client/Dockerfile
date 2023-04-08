FROM node:latest
WORKDIR /src/app
VOLUME ["/src/app/data"]
COPY . . 
EXPOSE 3000
RUN npm install
CMD ["npm", "run", "start"]