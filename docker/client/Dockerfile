FROM node:alpine as runtime

WORKDIR /app

# install package manager
RUN npm install -g pnpm


FROM node:alpine AS builder

WORKDIR /app

COPY client/package*.json ./
COPY client/pnpm-lock.yaml ./
COPY client/tsconfig.json ./
COPY client/svelte.config.js ./
COPY client/vite.config.ts ./

COPY client/src src
COPY client/static static

# install package manager
RUN npm install -g pnpm

RUN pnpm install
RUN pnpm run build

FROM runtime

WORKDIR /app

ARG timezone=Asia/Bangkok
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

COPY --from=builder /app/build build
COPY --from=builder /app/node_modules node_modules
COPY --from=builder /app/package*.json .

EXPOSE 3000

CMD [ "node", "build" ]