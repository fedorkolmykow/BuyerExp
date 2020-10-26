# Тестовое задание Авито в юнит buyer-experience

Для запуска сервиса:

`docker-compose up`

Пример запроса подписки на изменение цены объявления:

`
curl -d '{"mail":"powefes484@sekris.com","url":"https://www.avito.ru/moskva/avtomobili/volkswagen_amarok_2013_2034118645"}' -H "Content-Type: application/json" -X POST http://localhost:9000/subscribe
`

