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