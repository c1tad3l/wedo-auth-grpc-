<h1>Для работы с сервисом </h1>


```shell
docker pull postgres
```

```shell
docker volume create postgres_data
```

```shell

docker run --name postgres_container -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 -v postgres_data:/var/lib/postgresql/data postgres
```

<p> для произведения миграций </p>

```shell
task migrate 
```
<p> для очистки данных в базе  </p>

```shell
task table_down
```

<h1>
Example for request body data
</h1>

<p>Registration</p>

```json
{
    "email": "",
    "password":"",
    "phone": "",
    "dateOfBirth": "",
    "username": ""
}
```

<p>Login</p>


```json
{
    "email": "",
    "password": ""
}
```

<p>Check if is admin</p>

```json
{
    "user_uuid": ""
}
```

<p>To build dockerfile</p>

```shell
docker build -t  {Название аккаунта/ название сервиса: версия} -f Dockerfile .
```