# Bank

В данном проекте реализована работа со счетами клентов банка. 
Информация о счетах хранится в БД PostgreSQL Docker контейнерах через docker-compose файл

Для каждого счета:
- walletID - uuid счета 
- amount - баланс на счету 

## Makefile

В Makefile есть 3 цели:
- make all - запуск run_server
- make run_server - запуск сервера для работы со счетами (server) и БД (база PostgreSql)
- make test - unit-тесты, запускаются после run_server


## Конфигурация сервера и БД

Конфигурация находится в файле config.env в корне проекта.

При запуске, конфигурация автоматически считывается, поэтому их легко поменять при запуске.

В файле есть следующие параметры:
- SERVER_ADDRESS=адрес на котором работает сервер
- POSTGRES_HOST=хост базы данных
- POSTGRES_PORT= порт базы данных
- POSTGRES_USERNAME=имя пользователя базы данных
- POSTGRES_PASSWORD=пароль базы данных
- POSTGRES_NAME=имя базы данных

#### API

Server реализован на REST API. 

Реализованы слудующие методы:

- GET /api/v1/wallets/{WALLET_UUID} - метод получения информации о счетах по uuid.
После получения запроса, сервер делает запрос на БД о наличии и состоянии счета.
  Метод возвращает:
  - StatusCode 200 и информацию о счете.
  - StatusCode 400 и информацию об ошибке.
  - StatusCode 500 если произошла внутренняя ошибка на сервере

- POST /api/v1/wallet - метод изменения счета, который на входе получает UUID, тип операции и сумму, на которую надо изменить счет в зависимости от операции.
  Тип операции:
  - DEPOSIT -пополнение счета
  - WITHDRAW - снятие со счета
  
  Метод возвращает:
  - StatusCode 200 и информацию о результате выполнения.
  - StatusCode 400 и информацию об ошибке.
  - StatusCode 500 если произошла внутренняя ошибка на сервере

Остановить сервер можно сочетанием клавиш CRRL+C. После этого будет произведен Graceful Shutdown сервера c завершением всех запущенных goroutine, необходимых для более быстрой работы и решения проблемы конкурентных запросов.

Для удобной проверки работы REST API реализована документация Swagger по адресу /api/v1/docs/index.html


#### PostgreSQL

При запуске сервера, автоматически запускаются релеционные базы данных PostgreSQL в Docker. 
Dockerfile распологается в ./internal/database/
Настройка запуска баз реализована в файле config.env. 

В базах данных хранится информация о счетах

## Тестирование

В проекте реализовано автоматическое тестирование работы сервиса методом черной коробки.

Для запуска тестирования необходимо запустить Docker контейнеры сервера и БД make run_server и выполнить make test.

После тестирования можно ознакомиться с запущенными тестами, их результатом, ожидаемым результатом теста и полученным результатом. На сервере можно ознакомиться с логами в ходе выполенния запросов для тестирования