# ---- Base Node ----
FROM node:carbon AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
COPY package*.json ./
RUN npm ci

# ---- Copy Files/Build ----
FROM dependencies AS build
WORKDIR /app
COPY . /app
ARG EOS_ORIGIN
ENV EOS_ORIGIN=${EOS_ORIGIN}
ARG API_ORIGIN
ENV API_ORIGIN=${API_ORIGIN}
ARG SCRIPTS_ORIGIN
ENV SCRIPTS_ORIGIN=${SCRIPTS_ORIGIN}
ENV NODE_ENV=production
RUN npm run build
RUN npm run build-monaco
RUN npm run generate-docs

# --- Release with Alpine ----
FROM node:11.7-alpine AS release
WORKDIR /app
COPY --from=dependencies /app/package*.json ./
RUN npm ci --only=production
RUN mkdir public
COPY api api
COPY --from=build /app/dist ./public
COPY --from=build /app/static/monaco ./public/monaco
COPY --from=build /app/static/docs ./public/docs
COPY --from=build /app/server.js ./
EXPOSE 80
CMD ["npm", "start"]
