# Знай край
Backend для проекта Руси Сидящей на хакатоне Новой Газеты (https://projector2020.te-st.ru/)
<img src="https://github.com/semyon-dev/znai-krai/blob/master/img.png" alt="drawing" width="500"/>

### Используемые технологии на backend
Go 1.14, Gin, Google Maps API, Yandex Search API, Google Sheets API и другие

### Ссылки
* Наш сайт https://znaikrai.herokuapp.com/
* API endpoint https://api-znaikrai.herokuapp.com/
* Исходный код сайта: https://github.com/kniazevgeny/znaikrai

### Как запустить?
`go run main.go`

### Как скомпилировать в бинарник?
`go build main.go`

### Методы для клиентов

* метод для получения всех ФСИН учреждений \
`GET /places`

* BETA: отзывы с Google Maps \
`GET /reviews/:name`

* получение всех вопросов для создания новых нарушений со стороны клиента `(/form)` \
`GET /formQuestions`

* метод для создания новых нарушений (форм - заявок) \
`POST /form`

### TODO

- [ ] Метод для получения нарушений (форм)
- [ ] Аналитика
- [ ] Данные о коронавирусе в учреждениях
- [x] получение ФСИН учреждений

### License
znai-krai is licensed under the Creative Commons Attribution NonCommercial ShareAlike (CC-NC-SA)
