База данных Postgres по адресу: 127.0.0.1:5432/links , таблица t_links2

Файл с исходником: 'go_017.go'

1. Создание короткой ссылки POST /links
    Request { "url": "https://example.com/some/very/long/url" }
    
    Response { "short_code": "abc123" }
    
    Сервис должен:
    
    — сгенерировать short_code
    — сохранить ссылку в базу данных
    — вернуть short_code в ответе
    
    Сервис запускается по адресу: 127.0.0.1:443


2. Получение оригинальной ссылки

    GET /links/{short_code}
    
    Response
    
    {
      "url": "https://example.com/some/very/long/url",
      "visits": 15
    }
    
    При каждом запросе необходимо увеличивать счётчик visits.
