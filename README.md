# backend-trainee

Тестовое задание для стажировки в avito-tech 2023, направление backend.

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
