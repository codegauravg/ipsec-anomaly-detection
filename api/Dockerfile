# Multi-stage build for optimal image size and security
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .

# Production runner image
FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production

# Run as non-root user for security compliance
USER node

COPY --chown=node:node --from=builder /app/package*.json ./
COPY --chown=node:node --from=builder /app/node_modules ./node_modules
COPY --chown=node:node --from=builder /app/src ./src

EXPOSE 3000

CMD ["node", "src/server.js"]