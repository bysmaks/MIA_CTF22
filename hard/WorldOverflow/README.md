# WorldOverflow
### Сложность
Средне.
### Описание
Нашел интересный сервис, который исполняет желания! Можешь что-нибудь пожелать.

Коннект: nc ADDRESS 8888
### Решение
Уязвимость - переполнение буфера. Отправить 63 байта, которые переполнят буфер, перезатрут нулевой байт в конце. Выведется исходная строка, но уже вместе с флагом в конце.
### Флаг
CTF{y0ur_m4in_w15h}
### Деплой
Запустить скрипт ./start.sh, который поднимет докер.
