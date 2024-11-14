# Astral

### Get All Actice Functions
```http
GET http://localhost:8080/function
```


### Add a Functions
```http
POST http://localhost:8080/function
{
    "image": "nilspolek/echo-server",
    "port": 8080,
    "route": "/echo"
}
```

### Delete a Functions
```http
DELETE http://localhost:8080/function
# id = uuid of function
{
    "id": "e6268aa5-d312-46e1-a3ab-4ff928f5ab54"
}
```
