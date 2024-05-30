FROM node:alpine as runtime

WORKDIR /app

RUN npm install -g pnpm


FROM node:alpine AS builder

WORKDIR /app

COPY . .

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

COPY --from=builder /app/static static
COPY --from=builder /app/.svelte-kit .svelte-kit
COPY --from=builder /app/build build
COPY --from=builder /app/node_modules node_modules
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/pnpm-lock.yaml ./

CMD [ "node", "build" ]