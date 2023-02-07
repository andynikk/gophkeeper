1. Сервер запускается с флагами
-a адрес сервера
-d строка соединения с базой

Пример:
go run main.go -a localhost:8050 -d postgresql://postgres:postgres@localhost:5432/yapracticum

или параметры сеанса:
ADDRESS
CRYPTO_KEY

2. Клиент запускается с флагами
-a адрес сервера
-c файл с криптоключем

Пример:
go run main.go -a localhost:8080 -c e:\Bases\key\gophkeeper.xor

или параметры сеанса:
ADDRESS
DATABASE_URI
