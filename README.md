# USER MANAGER gRPC SERVER

In addition to the implementation of the necessary functionalities, the purpose of the service is to show how to organize a project and prevent it from turning into spaghetti code and how not to lose control when a microservice grows. The project is built according to the Onion architecture where the following layers can be noticed:

1. UI Layer (the most outer layer, represents the gRPC server with proto definitions and validations)
2. Business Logic Layer (services)
3. Domain Layer (all domain objects are at this layer)
4. Infrastructure Layer (provides functionality for accessing external systems like DB server, RMQ, Redis, etc. In our case, here is postgres integration with database migrations, repositories, RMQ integration and notification service responsible for notifying other services about user changes)

->app <br />
&nbsp;&nbsp;-> domain <br />
&nbsp;&nbsp;-> infrastructure <br />
&nbsp;&nbsp;-> services <br />
&nbsp;&nbsp;-> ui <br />

You can also notice 'app/config'. This is where environment variables are loaded immediately after starting the service. Because we can potentially use the environment variables in all layers, I left this aside, so it doesn't belong to any layer.

With this kind of architecture we get many advantages. Layers are not tightly coupled. It provides better maintainability as all the code depends on deeper layers or the centre. Improves overall code testability as unit tests can be created for separate layers without impacting other modules.

All of the methods in this project respect single responsibility principles. This is a very simple code where we can't really see that these rules are followed, but it's very important to follow them, because the code is clean and understandable, maintainable, reusable and easy to test.

You'll notice in the code that there is a lot of use of interfaces. Interface gives us a powerfull way to use abstraction. Most of the methods in this service are the implementation of some interface. I like to code this way because it will give me a nice and clean way to mock and write unit tests for any of these methods later. In several places you will notice the following structure:

-> repo <br />
&nbsp;&nbsp;&nbsp;-> mocks <br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- userRepoMock.go   (3) <br />
&nbsp;&nbsp;&nbsp;- userRepo_tests.go    (2) <br />
&nbsp;&nbsp;&nbsp;- userRepo.go          (1) <br />

  
userRepo.go (1) - the interface declaration and the actual methods implementation <br />
userRepo_tests.go (2) - unit tests for user repository <br />
userRepoMock.go (3) - mocked methods for userRepo <br />

Establishing good paradigms and consistent, accessible standards for writing clean code can help prevent developers from wasting many meaningless hours on trying to understand others (or their own) work.

# Database
For data storage is used Postgres server. Immediately after starting the service, a connection to the Postgres server is opened, a database is created (if it doesn't exist) and migrations are performed. I used Gorm ORM library for manipulation over the database. This library, built on the 'database/sql' package, is developer-friendly, easy-understandable and feature-rich. User is stored using required schema. Password is hashed. Nickname and email are unique. And the country code is composed of two letters.

# Notification system
In order to notify other services about changes to users, we use RabbitMQ open source message broker. The notification event is small and concise as it only contains a reference to the state that was changed - in our case user ID. Then consumers will determine if the change is relevant for them, and send request for the user. It uses a publish/subscribe mechanism, that represents an event-driven architecture, where any message published to a topic is immediately received by all of the subscribers to the topic. Go channel is used to pass the message from the NotificationService to the process responsible for publishing the messages to queue.

# Logging
For structured logging is used Zerolog library. Fast and simple logger dedicated to JSON output with stunning performance, avoiding allocations and reflection.

# Environment variables
Environment variables are defined in .env file. These env varibles are default ones and for development purposes, and can be overrided in docker-compose file. Godotenv library is used for loading and manipulating environment variables.

# Test coverage
To run the tests, you need to execute the following command from the root of the project:
    
    go test ./...

Or if you want to get verbose output that lists all of the tests and their results:
    
    go test -v ./..

Around 85% of the code is covered with unit tests. You can confirm that by runing:

    go test -v ./... -cover

Unit tests are written for critical and high-priority functions, functions that customers use and code that is frequently accessed. What is not covered is the database connection and the rabbit connection. This falls under system configuration so we need to have both servers (Posrgres and RMQ) active in order to write tests for this as well. In production code, connections to the database and external servers would generally be in a separate module, so that they can be reused and integrated into multiple services. It requires slightly different and more intensive testing than classic unit tests.

# Instructions on how to run:
Using a docker-compose file, you can easily run user manager in a container together with Postgres And RMQ containers on the same network. Dockerfile will copy all required files, build the service, run the tests and if everything goes well it creates a docker image. If you have installed docker engine on your computer, just run the following command:

    docker-compose up -d

To make requests, use some UI or tool for querying GRPC services. BloomRPC is a really good and simple tool. In UI layer there are 2 proto files. User proto and Health proto. Just import protos and examples of the requests will be created. But in case you use some others, I'll provide example requests:

1. Create user:

```json
{
  "firstname": "Aleksa",
  "lastname": "Vasiljevic",
  "nickname": "Ale94",
  "password": "pass",
  "email": "aleksa@gmail.com",
  "country": "RS"
}
```

2. Update user:

```json
{
  "id": "9eb24004-d476-4389-8a94-6e736aeb8011",
  "firstname": "Aleksa",
  "lastname": "Vasiljevic",
  "nickname": "Ale94",
  "password": "pass",
  "email": "aleksa@gmail.com",
  "country": "RS"
}
```

3. Get user page. You can filter by 3 possible values (country, createdFrom, creeatedTo):

```json
{
  "offset": 0,
  "limit": 10,
  "filter": {
    "country": "RS",
    "CreatedFrom": {
      "seconds": 1674590638
    },
    "CreatedTo": {
      "seconds": 1674590659
    }
  }
}
```

4. Delete user

```json
{
  "id": "027d0b3a-6053-476c-9bf4-c6494aac4df1"
}
```

As we agreed, really simple health check is implemented. A client can query the server’s health status by calling the Check method. A client can call the Watch method (which is not implemented) to perform a streaming health-check. The server will immediately send back a message indicating the current serving status. It will then subsequently send a new message whenever the service's serving status changes. To call Check method only provide:

```json
{
  "service": "name-of-the-service-that-you-want-to-check"
}
```

or in our case (because this is simple health check) you can provide empty service name:

```json
{
  "service": ""
}
```

# Possible extensions or improvements to the service for production
1. Database reconnection strategy in case of failure
2. Database migrations (really simple at the moment)
3. Database isolation level (in case multiple user manager instances are trying to update user)
3. RMQ setup
4. RMQ reconnection strategy
5. Better health checks (explained above)
6. Authentication/Authorization
7. CI/CD automated process
8. Logging unique value per request/flow (like correlation ID) so we can easily track errors in production
9. Logging level configurable via env variable
10. Integration of some kind of monitoring system (like Prometheus)

There would probably be more upgrades for production, but this is something that came to mind at the moment. :)

