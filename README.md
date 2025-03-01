# user-bookings-service
> This is a RESTful API created using Golang, PostgreSQL and Docker

### Requirements:
#### With Docker:
 ![docker](https://badgen.net/static/docker/@latest/purple)<br/>
 You can install Docker <a href="https://docs.docker.com/engine/install/">there</a>

#### Without Docker:
 ![golang](https://badgen.net/static/go/1.13/green?icon=github) ![postgresql](https://badgen.net/static/postgresql/@latest/)<br/>
 You can install Golang <a href="https://go.dev/doc/install">there</a><br/>
 You can install PostgreSQL <a href="https://www.postgresql.org/download/">there</a>

### Installing:
1. Clone repository 
2. In main directory:<br/>
   With Docker:
    for Windows users:
      ```bash
      docker-compose build
      docker-compose up
      ```
    for Linux users:
      ```bash
      docker compose build
      docker compose up
      ```
   Without Docker:
    Set PostgreSQL in pgAdmin (see in env file (.env))
    and <br/>
    
    ```
    go build
    
    //for windows
    bookingservice.exe

    //for linux
    ./bookingservice
    ```

### Entities:
 - **User (example)**:
```
{
  "id": 906,
  "username": "Andrew Tate",
  "password": "$2a$14$kv/sGmTWIlNYocbZqd88GuRsrOtKrs9bBFMM7N7HRNZ.qPxF.b.GG", //bcrypt hash
  "created_at": "2023-09-24 17:13:42",
  "updated_at": "2023-09-27 11:10:23"
}
```
 - **Booking (example)**:
```
{
  "id": 1021,
  "user_id": 906,
  "end_time": "2025-03-01 14:00:00",
  "start_time": "2025-03-01 20:42:00",
  "comment": "GOOOOOOOOLang"
}
```

### Requests
- /user/{id} [get]
  <br/>Get User by id
- /user [post]
  <br/>Create User from postForm: username, password
- /user/{id} [delete]
  <br/>Delete User and his bookings
- /user/{id} [put]
  <br/>Update (optional: username, password) User data by id (set new timestamp in update_at)

- /booking [post]
  <br/>Create User from postForm: user_id, start_time, end_time
- /booking [get]
  <br/>Get all bookings ordered by id
- /booking/{id} [get]
  <br/>Get Booking by id
- /booking/{id} [put]
  <br/>Update (optional: text, start_time, end_time) Booking data by id (set new timestamp in update_at)
- /booking/{id} [delete]
  <br/>Delete Booking by id