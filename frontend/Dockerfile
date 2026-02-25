# Build stage
FROM node:25-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source code and build
COPY . .
RUN npm run build

# Production stage
FROM node:25-alpine AS runner

WORKDIR /app

# Install all dependencies (needed for next.config.ts transpilation)
COPY package*.json ./
RUN npm ci

ENV NODE_ENV=production

# Copy built application from builder
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/next.config.ts ./next.config.ts

# node:alpineイメージにはnodeユーザーが存在するので、所有権を変更
RUN chown -R node:node /app

# 非rootユーザーに切り替え
USER node

EXPOSE 3000

CMD ["npm", "start"]
