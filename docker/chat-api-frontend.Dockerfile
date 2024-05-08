FROM node:18-alpine as node_base

WORKDIR /frontend

COPY frontend/ .
RUN rm -r node_modules/ .cache/ public/build/

RUN yarn install             
RUN yarn build          


FROM node_base as build_base

WORKDIR /frontend

COPY --from=node_base frontend/build/ ./build
EXPOSE 3000
ENV PORT 3000

CMD ["yarn", "start"]
