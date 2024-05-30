FROM node:18.12-alpine AS development

WORKDIR /app

COPY ./client .

RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build

CMD ["pnpm", "dev"]

FROM node:18.12-alpine

ARG timezone=Asia/Bangkok
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

COPY --from=development /app/build build

CMD [ "node", "build" ]