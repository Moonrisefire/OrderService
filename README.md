# OrderServiceToCheck
This version application "OrderService" to WBTECH to check.

1) Запуск через терминал docker-compose up -d
2) Подождать секунд 15.
3) Запуск через терминал go run cmd/main.go
4) Запуск через терминал go run cmd/send_order.go

Если нужно отправлять в БД другие данные либо заменить файл model.json, либо в файле send_order.go в 20 строке заменить название файла и заново запустить через терминал go run cmd/send_order.go

P.S. без предварительно запущенных docker-сервисов и основного файла не отправлять, тк сервер запускается в основном файле.


https://github.com/user-attachments/assets/de465034-8a94-422f-a909-f543828cb6e3

