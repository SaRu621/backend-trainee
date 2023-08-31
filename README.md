# backend-trainee

Тестовое задание для стажировки в avito-tech 2023, направление backend. 

Приложение работает на порту 9000.

В файле **config.go** содержатся настройки роутера и базы данных. Параметр Users нужен для автоматического добавления пользователей в сегмент (доп. задание 3).

Также, для подключения go-приложения к БД был добавлен скрипт **wait-for-postgres.sh**, который предпринимает попытку подключения до тех пор, пока этого не произойдет. Без этого скрипта инициализация БД может произойти позже, чем go-приложение попытается подключиться к ней, что приведет к отказу в подключении.

##  Как использовать приложение

### Запуск приложения
```console
git clone https://github.com/SaRu621/backend-trainee.git
cd backend-trainee
cd Avito
docker-compose up --build
```
### Получение сегментов пользователя:

```console
curl -X "GET" "http://localhost:9000/" \
-d $'{
    "id":"501"                                       
}'
```
### Добавление сегмента и некоторого процента пользователей в него:

```console
curl -X "POST" "http://localhost:9000/" -d $'{
    "slug":"AVITO_SEG_1",
    "percent":15
}'
```

### Добавление в сегменты и удаление из них пользователя:

```console
curl -X "PUT" "http://localhost:9000/" -d $'{
    "add":["AVITO_SEG_1", "AVITO_SEG_2"],
    "del":["AVITO_SEG_3"],
    "id":"501",
    "exp":[{
            "y":"1",
            "mo":"2",
            "d":"3",
            "h":"4",
            "mi":"5"
           },
           {
            "y":"",
            "mo":"",
            "d":"",
            "h":"",
            "mi":""
           }            
    ]
}'
```
**Пояснение**: в данном случае пользователь с id 501 будет:

- Добавлен в сегмент AVITO_SEG_1 на 1 год, 2 месяца, 3 дня, 4 часа и 5 минут (реализация TTL);  

- Добавлен в сегмент AVITO_SEG_2 без ограничений по времени;  

- Удален из сегмента AVITO_SEG_3.

### Удаление сегмента:
```console
curl -X "DELETE" "http://localhost:9000/" \
-d $'{
    "slug":"AVITO_SEG_4"                                       
}'
```
### Запрос URL-ссылки на csv-таблицу событий, происходящих с пользователем за определенный год-месяц:
```console
curl -X "GET" "http://localhost:9000/report/" \
-d $'{
    "id":"501",
    "year":"2023",
    "month":"8"         
}'
```

### Получение csv-таблицы:

```console
curl -X "GET" "http://localhost:9000/report/"
```
Пример полученной таблицы:

```console
[
    "501,AVITO_SEG_1,added,2023-08-30 09:49:06 ",
    "501,AVITO_SEG_2,added,2023-08-30 13:21:15 ",
    "501,AVITO_SEG_3,added,2023-08-30 13:41:33 ",
    "501,AVITO_SEG_3,deleted,2023-08-30 13:42:16 "
]
```
### Удаление всех csv-таблиц:

```console
curl -X "DELETE" "http://localhost:9000/report/"
```
## Зависимости

- [gin](https://github.com/gin-gonic/gin) - HTTP web фреймворк
  
- [pq](https://github.com/lib/pq) - Postgres драйвер

## Инструменты
- [Golang](https://go.dev/) - язык программирования
  
- [Postgres](https://www.postgresql.org/) - база данных

- [docker](https://www.docker.com/) - технология виртуализации
