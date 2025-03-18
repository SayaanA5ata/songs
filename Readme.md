# Test task 

### Technologies
- Chi
- Gorm
- PostgreSQL
- Docker
- docker-compose
- Swagger

### Запуск БД
- БД можно запустить командой 
```docker-compose up```
- Миграции нужно накатить отдельно, это автомиграция 
### Миграция 
```go run migrations/auto/auto.go```

- Переменные для сервера и БД находятся в файле .env

<h2>Структура проекта</h2>

<code>
/song
├───cmd - пакет с точкой входа
├───configs - пакет с конфигами
├───docs - сгенерированный сваггер
├───internal
│   ├───domain - доменный слой
│   ├───interfaces - слой интерфейсов
│   ├───handlers - слой обработчиков
│   ├───server - логика сервера
│   └───service - слой бизнес-логики
│───migrations - пакет миграций
│    └───auto - автомиграция
└───pkg - модули, которые можно переиспользовать
    ├───db - модуль БД
    ├───dsn - получение DSN, спорная штука
    ├───middleware - пакеты middleware
    ├───req - модуль request
    └───res - модуль response
</code>