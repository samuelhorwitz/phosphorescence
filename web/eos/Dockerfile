# ---- Base Node ----
FROM node:carbon AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
WORKDIR /app/eos
COPY ./eos/package*.json ./
RUN npm ci

# ---- Copy Files/Build ----
FROM dependencies AS build
WORKDIR /app
COPY ./common /app/common
COPY ./eos /app/eos
WORKDIR /app/eos
ARG PHOSPHOR_ORIGIN
ENV PHOSPHOR_ORIGIN=${PHOSPHOR_ORIGIN}
ENV NODE_ENV=production
RUN npm run build

# --- Release with Alpine ----
FROM node:11.7-alpine AS release
WORKDIR /app
COPY --from=dependencies /app/eos/package*.json ./
RUN npm ci --only=production
COPY --from=build /app/eos/dist ./public
COPY --from=build /app/eos/server.js ./
EXPOSE 80
CMD ["npm", "start"]
