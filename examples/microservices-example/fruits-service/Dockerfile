# Use node 4.4.5 LTS
FROM node:8-alpine
#ENV PORT 8000

# Copy source code
COPY . /app

# Change working directory
WORKDIR /app

# Install dependencies
RUN npm install

# Expose API port to the outside
EXPOSE 9999

# Launch application
CMD ["npm","start"]