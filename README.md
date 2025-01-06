# Задание
Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе. Модель данных можно найти в конца задания.

Что нужно сделать:

<ol>
  <li>Развернуть локально PostgreSQL</li>
  <ol>
    <li>Создать свою БД</li>
    <li>Настроить своего пользователя</li>
    <li>Создать таблицы для хранения полученных данных</li>
  </ol>
  <li>Разработать сервис</li>
  <ol>
    <li>Реализовать подключение к брокерам и подписку на топик orders в Kafka</li>
    <li>Полученные данные записывать в БД</li>
    <li>Реализовать кэширование полученных данных в сервисе (сохранять in memory)</li>
    <li>В случае прекращения работы сервиса необходимо восстанавливать кэш из БД</li>
    <li>Запустить http-сервер и выдавать данные по id из кэша</li>
  </ol>
  <li>Разработать простейший интерфейс отображения полученных данных по id заказа</li>
</ol>
