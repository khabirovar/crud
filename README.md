CRUD
================================
Simple CRUD implementation task

Required list of api:
- Get item by id (GET /books:id)
- Get list of items (GET /books)
- Create item (POST /books)
- Update item (PUT /books)
- Delete item by id (DELETE /books:id)

Other requirements:
- PostgreSQL as data storage
- Using JSON for data transfer
- Dockerize 
- Docker compose
- List of requests for testing
- Caching requests with in memory cache


Simple db with one table:
```sql
create table books (
    id serial not null unique,
    title varchar(255) not null unique,
    author varchar(255) not null,
    publish_date timestamp not null default now(),
    rating int not null
);
```

For run application use docker compose:
```shell
docker-compose up -d
```

