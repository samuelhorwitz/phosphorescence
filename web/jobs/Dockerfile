FROM node:11.7-alpine AS release
WORKDIR /app
COPY ./jobs ./
RUN npm ci --only=production
ENV NODE_ENV=production
CMD ["npm", "start"]
