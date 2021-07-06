# balance-service
***
## HTTP Server:
### Ручки:
+ Проверка соединения сервиса с базой данных и первичной синхронизации с линиями:  
  Request: **[GET] /ready**

Response:
<pre>
200

500
</pre>

## GRPC Server:
**_api/kiddy.proto_**
### Пример:
**Request**
<pre>
  call SubscribeOnSportsLines
  [repeated] lines (TYPE_STRING) => soccer
  [repeated] lines (TYPE_STRING) => baseball
  time (TYPE_DOUBLE) => 5
</pre>
**Response**
<pre>
  {
    "linesDeltas": {
      "soccer": 1.97,
      "baseball": 2.024
    }
  }
//wait 5 seconds
  {
    "linesDeltas": {
      "soccer": -0.702,
      "baseball": 0.024
    }
  }
</pre>
## Переменные окружения:

    * HTTP_PORT - порт HTTP сервера (default="8080")
    * GRPC_PORT - порт GRPC сервера (default="8081")
    * LP_ADDRESS - адрес lines provider (default="http://localhost:8000")
    * STORAGE - адрес redis хранилища (default="localhost:6379")
    * STORAGE_PASSWORD - пароль redis хранилища (default="")
    * B_TIME - интервал обновления коэффициента линии baseball (default="1")
    * F_TIME - интервал обновления коэффициента линии football (default="1")
    * S_TIME - интервал обновления коэффициента линии soccer (default="1")
***
## Запуск:

>docker-compose up
