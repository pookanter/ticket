FROM node:20-alpine AS development

WORKDIR /app

COPY ./client .

RUN npm install -g pnpm
RUN pnpm install
RUN npm uninstall @rollup/rollup-linux-arm64-musl
RUN npm i
RUN pnpm build

CMD ["pnpm", "dev"]

FROM node:20-alpine AS production

ARG timezone=Asia/Bangkok
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

COPY --from=development /app/build build
COPY --from=development /app/package.json .

CMD [ "node", "build" ]