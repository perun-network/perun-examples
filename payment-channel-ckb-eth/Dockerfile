FROM node:20-alpine

WORKDIR /hardhat

COPY package.json ./
COPY hardhat.config.js ./

RUN npm install

EXPOSE 8545

CMD ["npx", "hardhat", "node", "--hostname", "0.0.0.0"]
