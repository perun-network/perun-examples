FROM node:20-alpine

# Create app directory
WORKDIR /hardhat

# Install required dependencies
COPY ../package.json ../yarn.lock ./
RUN yarn install --frozen-lockfile

# Copy Hardhat config and scripts
COPY ../hardhat.config.js ./

# Set up environment
ENV HOST=0.0.0.0
ENV PORT=8545

EXPOSE 8545

ENTRYPOINT ["sh", "-c", "npx hardhat node --hostname 0.0.0.0 --port 8545"]