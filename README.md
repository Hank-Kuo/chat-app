# Chat app
## Introduction
Design a high performance and scalable chat-app in distribution system. 

Here is two arch I came up with, 
1. Using grpc as communication channel to send message directly
- This method can't ensure that the message send to users, it may miss some messages.   

![chat-app-arch-1](doc/chat-app-arch-1.png) 
![chat-app-websocket-1](doc/chat-app-websocket-1.png) 

2. Using message queue to broadcast message to each user.
- Each user can subscribe the channel they joined.

![chat-app-arch-2](doc/chat-app-arch-2.png) 


For this project, I adopt arch 2, and I use channel_id as each topic.



### How to run

In developer mod: 
```bash
make docker-up
make run-api

cd frontend && yarn dev
```

In docker mod: 
```bash
make run-docker

```

## Frontend Library

- react: frontend framework
- remix: react SSR lib
- emotion: js-based CSS lib in SSR

## Backend Library

- kafka: wurstmeister/kafka
- zookeeper: wurstmeister/zookeeper
- kowl: UI dashboard for kafka
- cassandra: NoSQL database
- redis: NoSQL database
- prometheus: monitor & alert tool
- node_exporter
- grafana: UI dashboard
- jaeger: trace serve request
- gin: go http server
- grpc: go grpc server
- gorilla: go websocket server
- go-kakfa: go kafka's third-party lib
- go-migrate: db migration lib
- sqlx: go db third-party lib
- viper: go config lib

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
  - DI concept: high module shouldn't depend on lower module (IoC is deisng principle, the IoC core concept is try to transfer controller into container or framework, the DI is one of implement method)
    - Using interface(Factory pattern) to abstract the modules, then inject interface into dependency module, it help us easy to maintain our modules, and decouple all modules complexity. 


For addressing the massive message traffic, I adopt nosql database, cassandra. Cassandra can tolerate heavy write traffic, it's a scable, fault-tolerant and performance distributed database, using CAP theorem to satisfy eventually consistent. The CQL is how we query cassandra language. Why we don't use mongoDB, the mongoDB more prefer heavy read traffic system, but in chat-app scenario, it more care about dealing with heavy write, that's why I choose Cassandra. 
Here is useful how to design Cassandra schema document: https://tachunwu.github.io/posts/cassandra/

For dealing the massive read traffic, the redis is very good way to alleviate db access traffic.
- cache strategy 
  - cache aside 
    - check cache at first time, if data isn't in cache, then retrieved from database and stored in cache.
    - Good: read-heavy
  - read through 
    - cehck cache, if not hit, cache can rerieved from database and store in cache. (different with cache aside, retreived part, it's from cache) 
    - Good: read-heavy, first time may delay. 
  - write through 
    - write data in cache and in DB at the same time, it can ensure data in cache is consist as DB. (But cost extra write latency) 
    - Good: data consistency, but the latency will increase.  
  - write behind  (write back)
    - write data in cache at first, then write data in DB later (asynchronously updated)
      - To handling data missing problem, always use kafka together.
    - Good: write-heavy

For this system, I adopt cache aside and write throguh for `channel` cache to handle high traffic query. Also in this system is very good for using write behind for `message/reply`, beacause `message/reply` is write-heavy feature, I don't have free time to implement this part.

Useful cache strategy document: 
- https://medium.com/@mmoshikoo/cache-strategies-996e91c80303
- https://homuchen.com/posts/databse-chache-strategies/

FYI cache some issuse
- cache avalanche
- cache penetration 
- cache hotspot


Kafka is one of popular message queue, it's scalable, high performance, high throughput and fault-tolerant open source. For message queue, we can use queue to accomplish pub/sub system / streaming pipeline / event driven arch / decouple system complexity. The kafka is more suit for high throughput, it can deal with million message per second, it always uses in event log process, streaming/data pipeline, async job. the RabbitMQ is more suit for complexity routing system (routing and exchange), but it only can deal with throunds messages per second, it supports priority queue (FIFO), it always uses in event driven arch, decouple component, async job. 
The kafka has some compnents we need to know: producer, consumer, broker, topic, partition. 



Kafka 
- partition / replicas
- exactly once 
- 




For addressing hotspot problem for Cassandra, 

here an example, 
Sharding -> avoid hotspot problem 


primary key
message table: chaneel_id, bucket 
reply table: message_id 

cluster key
message table: message_id
reply table: reply_id

Secondary Index


Addressing large data in the database, for the message & reply `get` API, I adpot cursor to alleviate scan whole database. Following the cursor and pagination, the pagination has a big fault for database, the pagination scan the table, it means the query need to read whole table, it cost lots of performance for query. 


How to keep communicate with client side server, there are serval way to achieve, ex: long-polling (Comet), websocket,graphQL. THe long-polling is very inefficient way to communicate with server, it costs lots of connection, the server and client always build TCP/IP connect in every connection. The Websocket is bidirectional communication protocal, it's very suit for chat-app scenario, it keep the connection after they build the connection at first time. 



Snowflake


### API document 
Auth API
- login/register
  - login:
    - path: /api/login
    - request
    ```json
    { 
      "email": "test@gmail.com", 
      "password": "test123456"
    }
    ```
    - response
    ```json
      {
        "status":"success",
        "message":"login successfully",
        "data":{
          "id":"257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
          "name":"test",
          "email":"test@gmail.com","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImhhbmsiLCJFbWFpbCI6Imhhbmtfa3VvQGdtYWlsLmNvbSIsIlVzZXJJRCI6IjI1N2U0Y2FmLWZiNGItNDNhMS1hNGIzLWNjYTk0ZjU4M2JkNSIsImlzcyI6Imhhbmsta3VvIiwiZXhwIjoxNzE0MTM5NjUzLCJpYXQiOjE3MTQwNTMyNTN9.qvFYiQ994oyNp2E9oQpUnj8rden_bvr2o_9F_qXYkLs","created_at":"2023-12-01T05:52:17.581988Z"
        }
    }     
    ```
  - register:
    - path: api/register
    - request
    ```json
    {
      "username": "test", 
      "email": "test@gmail.com", 
      "password": "test123456"
    }
    ```

Channel API
- channel: get/getUserChannel/create/join
  - get
    - path: /channel
    - request
    ```json
    ```
    - response
    ```json
    ```
  - getUserChannel:
    - path: /user/channel
    - request 
    ```json
    ```
    - response
    ```json
    ```
  - create
    - path: /channel
    - body: name/channel_id
    - request
    ```json
    ```
    - response
    ```json
    ```
  - Join:
    - path: /channel/join
    - body: channel_id
    - request
    ```json
    ```
    - response
    ```json
    ```

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

  


### Database schema 
PostgreSQL tables
users table
| id            | username      |    email               | password      | created_at    | login_time    |
| ------------- | ------------- | ---------------------- | ------------- | ------------- | ------------- |
| uuid (PK)     | VARCHAR(255)  | VARCHAR(255) (unique)  | VARCHAR(255)  | TIMESTAMP     | TIMESTAMP     |

channel
| id            | name          | user_id   | created_at    | 
| ------------- | ------------- | --------- | ------------- | 
| uuid (PK)     | VARCHAR(255)  | uuid (FK) | TIMESTAMP     | 


user_to_channel
| user_id       | channel_id    | created_at    |
| ------------- | ------------- | ------------- |
| uuid (PK/FK)  | uuid (PK/FK)  | TIMESTAMP     |

Cassandra

message table
| user_id       | channel_id    | created_at    |
| ------------- | ------------- | ------------- |
| uuid (PK/FK)  | uuid (PK/FK)  | TIMESTAMP     |

reply table
| user_id       | channel_id    | created_at    |
| ------------- | ------------- | ------------- |
| uuid (PK/FK)  | uuid (PK/FK)  | TIMESTAMP     |

## Todo
- websocket error handler with kafka 
  - retry queue
- notification consumer 
  - send message to notification 
- write message & reply to redis
  - get message and reply from redis -> cursor
- docker-compose for all service 
- screen shots
Need to run stateful server
- saga pattern 