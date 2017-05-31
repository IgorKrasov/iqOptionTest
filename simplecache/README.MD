# simple in-memory key/value cache

Реализована поддержка 3х типов данных string, list, dict

##Общие доступные операции

### keys
 Вернет все ключи

 request:
 ```
curl -X GET \
  http://<host>/keys \
  -H 'cache-control: no-cache' \
 ```

 success response:
 ```
 //http.StatusCode: 200
 "test"
 ```
 failure response:
 ```
 http.StatusCode: 404, 400
 <error message>
 ```

### unset
  Удалит значение по ключу

  request:
  ```
  curl -X DELETE \
    http://<host>/unset/<key> \
    -H 'cache-control: no-cache' \
  ```

  response:
  ```
  //http.StatusCode: 201
  {
    "result": "success"
  }
  ```


## Cтроки
 Возможные операции set, get, unset
 Позволяют сохранить и получить данные любого типа в виде строки по ключу

###set

Сохранит значение value по ключу key, имеется ttl, которое задается в expired.
Если передать значение epired == 0, то ключ не будет "протухать"

request:
```
curl -X POST \
  http://<host>/set \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"key": "test:key",
	"expired": 1500,
	"value": "test"
}'
```
response:
```
//http.StatusCode: 201
{
  "result": "success"
}
```

### get
Позволяет получить значение по ключу key

request:
```
curl -X GET \
  http://<host>/get/<key> \
  -H 'cache-control: no-cache' \
```

success response:
```
//http.StatusCode: 200
"test"
```
failure response:
```
//http.StatusCode: 404, 400
<error message>
```

## Списки
    Реализована возможность сохранять по ключу список, добавлять в него значенияб просматривать и извлекать

### rpop

Если значения по переданному ключу не существует создаст новый список и запишет в него value.
Если значение по ключу уже существует и тип у него список, добавит значение value в конец списка
В противно случае будет ошибка

request:
```
curl -X POST \
  http://<host>/rpush \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"key": "testlist:key",
	"expired": 5,
	"value": "test_list123"
````

success response:
```
//http.StatusCode: 201
{
  "result": "success"
}
```

failure response:
```
//http.StatusCode: 404, 400
<error message>
```

###pop
Извлекает крайний значение из списка по ключу
В случае если тип данных по ключу не список - вернет ошибку
В случае если извлекается последний ключ из списка значение по ключу и ключ удаляются из хранилища

request:
```
curl -X GET \
  http://<host>/pop/<key> \
  -H 'cache-control: no-cache' \
```
success response:
```
//http.StatusCode: 200
"test"
```
failure response:
```
//http.StatusCode: 404, 400
<error message>
```

### lgetall
Вернет все значения из списка по заданному ключу
Если по ключу храниться не список вернет ошибку

request:
```
curl -X GET \
  http://<host>/lgetall/<key> \
  -H 'cache-control: no-cache' \
```
success response:
```
//http.StatusCode: 200
[
  "test_list1234",
  "test_list123"
]
```

failure response:
```
//http.StatusCode: 404, 400
<error message>
```

### lget
Вернет значение по ключу и номеру элемента в списке
Если по ключу храниться не список  или нет элемента по номеру вернет ошибку

request:
```
curl -X GET \
  http://<host>/lget/<key>/<i> \
  -H 'cache-control: no-cache' \
```
success response:
```
//http.StatusCode: 200
"test"
```
failure response:
```
//http.StatusCode: 404, 400
<error message>
```

### Словари

Позволяет по ключу сохранять словарь ключ/значение

##hset
Если значение по переданному ключу не существует создаст новый список и запишит в него данные
Если сушествует то, в случае отсутсвия в словаре пары ключ/значение добавит их в словарь,
при наличии - перезапишет значение в словаре

request:
```
curl -X POST \
  http://<host>/hset \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"key": "testH:key",
	"expired": 1500,
	"value": {
		"name": "Igor",
		"age": 28
	}
}'
```

success response:
```
//http.StatusCode: 201
{
  "result": "success"
}
```

failure response:
```
//http.StatusCode: 404, 400
<error message>
```

### hgetall
   Вернет весь словарь по ключу, если по ключу сохранен не словарь вернет ошибку

request:
```
curl -X GET \
  http://<host>/hgetall/<key> \
  -H 'cache-control: no-cache' \
```

success response:
```
//http.StatusCode: 201
{
   "name": "Igor",
   "age": 28
}
```

failure response:
```
http.StatusCode: 404, 400
<error message>
```

### hget
 Вернет значение из словаря по ключу
 В случае если значение в словаре по переданному ключу нет или тип данных по ключу не словарь - вернет ошибку

 request:
 ```
curl -X GET \
  http://<host>/hget/<key>/<dictKey> \
  -H 'cache-control: no-cache' \
 ```

success response:
```
//http.StatusCode: 201
"Igor"
```

failure response:
```
//http.StatusCode: 404, 400
<error message>
```