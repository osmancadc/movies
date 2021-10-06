# Movies

Dillinger this projects allows you to login/register using JWT and CRUD movies.

## Installation

Dillinger requires [golang](https://golang.org/dl/) v1.17.1+ to run.

Install the dependencies and devDependencies and start the server.

```sh
cd movies
go mod init serverMovies
go mod tidy
```
Run the server

```sh
go run .
```

## How to use

### Register user
```bash
    <host>:3000/users/register [POST]
```
```javascript
{ 
    "name": <string>
    "email": <string>
    "password": <string>
}
```
---
### Login user
```bash
    <host>:3000/users/login [POST]
```
```javascript
{ 
    "email": <string>
    "password": <string>
}
```

##### Login user will return a token, that token must be used in the Authorization header, in all of the following requests 
---
### Create a movie
```bash
    <host>:3000/movies/create [POST]
```
```javascript
{ 
    "name": <string>
    "duration": <integer>
    "gender": <string>
    "premiere_year": <string>
    "sales": <integer>
}
```
---
### Delete a movie
```bash
    <host>:3000/movies/delete [DELETE]
```
```javascript
{ 
    "id": <integer>
}
```
---
### Delete all movies
```bash
    <host>:3000/movies/delete/all [DELETE]
```
---
### Update sales of a movie
```bash
    <host>:3000/movies/update/sales [PATCH]
```
```javascript
{ 
    "id":<integer>
    "sales":<integer>
}
```
---
### Like a movie
```bash
    <host>:3000/movies/like [POST]
```
```javascript
{ 
    "id":<integer>
}
```
---
### Get public movies
```bash
    <host>:3000/movies/get/public [GET]
```
---
### Get private movies
```bash
    <host>:3000/movies/get/private [GET]
```
---
### Get liked movies
```bash
    <host>:3000/movies/get/liked [GET]
```
---
### Get random number
```bash
    <host>:3000/random [GET]
```




