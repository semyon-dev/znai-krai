# Знай край
Бэкенд проекта Руси Сидящей - Знай Край. Проект был начат на хакатоне Новой Газеты (https://projector2020.te-st.ru/)
<img src="https://github.com/semyon-dev/znai-krai/blob/master/img.png" alt="drawing" width="500"/>

### Используемые технологии на бэкенде
* Go 1.14
* MongoDB
* Gin
* Google Maps API
* Yandex Search API
* Google Sheets API

### Ссылки
* Сайт проекта https://znaikrai.herokuapp.com/
* API endpoint https://api-znaikrai.herokuapp.com/
* Исходный код сайта: https://github.com/kniazevgeny/znaikrai

### Как запустить?
`go run main.go`

### Как скомпилировать?
`go build main.go`

### API Методы

* методы для получения ФСИН учреждений
Все сразу `GET /places`

Конкретное `GET /places/<id>`
Нарушения есть только для конкретных учреждений

* получение всех ФСИН учреждений у которых есть информация по коронавирусу \
`GET /corona_places`

* получение всех нарушений сразу \
`GET /violations`

* BETA: отзывы с Google Maps \
`GET /reviews/<name>`

* аналитика по разным параметрам (общая статистика) \
`GET /analytics`

* получение всех вопросов для создания новых нарушений со стороны клиента `(/form)` \
`GET /formQuestions`

* метод для создания новых нарушений (форм - заявок) \
`POST /form`

### License
znai-krai is licensed under the Creative Commons Attribution NonCommercial ShareAlike (CC-NC-SA)
