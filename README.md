# Chat app
## Introduction
Design a high performance and scalable chat-app in distribution system. 


## Frontend Library

- react: frontend framework
- remix: react SSR lib
- emotion: js-based CSS lib in SSR

## Backend Library

kafka: wurstmeister/kafka
zookeeper: wurstmeister/zookeeper
kowl: UI dashboard for kafka
cassandra: NoSQL database
redis: NoSQL database
prometheus: monitor & alert tool
node_exporter
grafana: UI dashboard
jaeger: trace serve request
gin: go http server
grpc: go grpc server
go-kakfa: go kafka's third-party lib
go-migrate: db migration lib
sqlx: go db third-party lib
viper: go config lib

## Frontend
Using remix in typescript for SSR, the entrypoint of frontend is `root.tsx`, and `pages` are in routes folder. Here is Remix doc: https://remix.run/. 

Frontend structure: 
- apis
  - provide the apis connected with backend service
- commponents
  - sidebar
    - users' channels sidebar
- context
  - messageContext
    - connect websocket to send data or received data
- lib
  - hook: some hook function
  - utils: session, date, etc ...
- routes
  - _index
    - path: `/`
    - chat homepage
  - login
    - path: `/login`
    - login page
  - register
    - path: `/login`
    - register page


Remix a little knowledge:
- Loaders
  - SSR load data 
  - hook ref: useLoaderData
- Action
  - SSR post data
  - in Remix, we always use formData to get data which we want to pass by api
  - hook ref: useActionData
- Fetch
  - SSR post data
  - Different with action, the fetch work on when the url shouldn't change
  - refer: https://remix.run/docs/en/main/discussion/form-vs-fetcher
- Outlet
  - Renders the matching child route of a parent route.
  - we can pass conext through outlet, or use Context to cover outlet
- entry.client
  - using `ReactDOM.hydrateRoot` to hydrate the markup
- entry.server
  - render the markup for the current page using a <RemixServer> element with the context and url for the current request

## Backend
Using `gin` as http server router, `gorilla/websocket` as webscoekt server and grpc as internal communication server.

The Backend adopt DDD architecture using IoC and DI design pattern, you can refer `internal/api`
- DDD
  - delivery: provide different transport layer protocal 
  - repoitory: communicate with infra layer or data layer  
  - service: bussiness logic bewteen delivery and repoitory  
- IoC and DI
  - DI concept: high module shouldn't depend on lower module
    - Using interface(Factory pattern) to 
    - Factory

IoC 的範疇包含 DI，但不僅限於 DI。
decouple


For dealing the massive message traffic, I adopt nosql database, cassandra. Cassandra can 

For dealing the massive redis

Cache strategy
 
use cassandra
use websocket
use kakfa
use Snowflake ID

primary key
cluster key
Secondary Index


- cursor pagination
  - websocekt
  - queue

### API document 
Auth API
- auth: login/register
  - login:
    - path: /api/login
    - body: email/password
    - request
    ```json
    { 
      "email": "hank_kuo@gmail.com", 
      "password": "hank123456"
    }
    ```
    - response
    ```json
      {
        "status":"success",
        "message":"login successfully",
        "data":{
          "id":"257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
          "name":"hank",
          "email":"hank_kuo@gmail.com","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImhhbmsiLCJFbWFpbCI6Imhhbmtfa3VvQGdtYWlsLmNvbSIsIlVzZXJJRCI6IjI1N2U0Y2FmLWZiNGItNDNhMS1hNGIzLWNjYTk0ZjU4M2JkNSIsImlzcyI6Imhhbmsta3VvIiwiZXhwIjoxNzE0MTM5NjUzLCJpYXQiOjE3MTQwNTMyNTN9.qvFYiQ994oyNp2E9oQpUnj8rden_bvr2o_9F_qXYkLs","created_at":"2023-12-01T05:52:17.581988Z"
        }
    }     
    ```
  - register:
    - path: /register
    - body: email/username/password
    - request
    ```json
    {
      "username": "hank", 
      "email": "hank_kuo@gmail.com", 
      "password": "hank123456"
    }
    ```

Channel API
- channel: get/getUserChannel/create/join
  - get
    - path: /channel
  - getUserChannel:
    - path: /user/channel
  - create
    - path: /channel
    - body: name/channel_id
  - Join:
    - path: /channel/join
    - body: channel_id

Message API
- message:
  - send
    - path: /api/message/ws
    - request body: 
    ```json
    {
      "user_id": "257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
      "username": "test",
      "channel": "b37c4896-70e5-4a94-bbab-7de13e5e41ff",
      "content": "test message",
    }
    ```
    - response body:
    ```json
    {
      "message_id": 1783496863102533632,
      "user_id": "257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
      "username": "test",
      "channel": "b37c4896-70e5-4a94-bbab-7de13e5e41ff",
      "content": "test message",
      "created_at":"2024-04-25T14:02:42.590621Z"
    }
    ```

  - get
    - path: /message

- Reply:
  - send
    - path: /reply
  - get
    - path: /reply

  

## Todo
- in broadcast, call grpc api
- test all function
- consumer 
  - send message to notification 
- docker-compose for all service 
- screen shots