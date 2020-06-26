# Знай край
Бэкенд социального и некоммерческого проекта Руси Сидящей - Знай Край. Проект был начат [на хакатоне Новой Газеты](https://projector2020.te-st.ru/).

<img src="https://github.com/semyon-dev/znai-krai/blob/master/img.png" alt="drawing" width="500"/>

### Используемые технологии на бэкенде
* Go 1.14
* MongoDB
* Gin
* Google Sheets API
* Google Maps API
* Yandex Search API

### Ссылки
* Сайт проекта [znai-krai.zekovnet.ru](https://znai-krai.zekovnet.ru/) и [znaikrai.herokuapp.com](https://znaikrai.herokuapp.com/)
* Публичный API https://api-znaikrai.herokuapp.com/
* Исходный код сайта: https://github.com/kniazevgeny/znaikrai

### Лицензия
znai-krai is licensed under the Creative Commons Attribution NonCommercial ShareAlike (CC-NC-SA)

Лицензия позволяет другим перерабатывать, исправлять и развивать проект на некоммерческой основе, до тех пор пока они упоминают оригинальное авторство и лицензируют производные работы на аналогичных лицензионных условиях. Все новые работы, основанные на этом проекте, должны иметь эту же лицензию, поэтому все производные работы также должны носить некоммерческий характер.

### Как запустить?

Необходимо добавить переменные окружения (смотри config)

`go run main.go`

Или скомпилировать в единый бинарник:

`go build main.go`

### Документация к API методам
##### Публичные методы

Протокол: HTTP, формат данных: JSON

* методы для получения ФСИН учреждений \
Все сразу `GET /places` \
Конкретное `GET /places/<id>`
Нарушения (в том числе по короне) есть только для конкретных учреждений

* получение всех нарушений у которых есть информация по коронавирусу \
`GET /corona_places`

* получение всех нарушений \
`GET /violations`

* аналитика по разным параметрам (общая статистика) \
`GET /analytics`

* пояснения по разным параметрам (скорее для аналитики) \
`GET /explanations`

* получение всех вопросов для создания новых нарушений со стороны клиента `(/form)` \
`GET /formQuestions`

##### Закрытые методы ([свяжитесь](https://t.me/semyon_dev), чтобы получить доступ)

* сообщение новых нарушений (форм, заявок)
```
POST /form
!!! place_id string !!!
Параметры нужно получать из GET /formQuestions
```

* создание сообщений по коронавирусу (форм - заявок)
```
POST /form_corona
Параметры:
name_of_fsin string (название МЛС)
place_id string
region string
info string
```

* сообщение ошибок/багов
```
POST /report
Параметры:
email string
bug string
place_id string
name_of_fsin string
```

* подписка на email рассылку
```
POST /mailing
Параметры:
name string (Имя)
email string (обязательный параметр)
```

* Deprecated: отзывы с Google Maps
```
GET /reviews/<name>
```