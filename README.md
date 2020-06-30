# Знай край
Бэкенд социального и некоммерческого проекта Руси Сидящей - Знай Край. Проект был начат [на хакатоне Новой Газеты](https://projector2020.te-st.ru/).

<img src="https://github.com/semyon-dev/znai-krai/blob/master/img.png" alt="drawing" width="500"/>

### Используемые технологии на бэкенде
* Golang v1.14
* MongoDB
* Gin
* Google Sheets API
* Google Maps API
* Yandex Search API

### Ссылки
* Сайт проекта [znai-krai.zekovnet.ru](https://znai-krai.zekovnet.ru/)
* Публичный API https://api.znai-krai.zekovnet.ru/ и https://api-znaikrai.herokuapp.com/
* Исходный код сайта: https://github.com/kniazevgeny/znaikrai

### Лицензия
znai-krai is licensed under the [Creative Commons Attribution NonCommercial ShareAlike (CC-NC-SA)](https://github.com/semyon-dev/znai-krai/blob/master/LICENSE)

Лицензия позволяет другим перерабатывать, исправлять и развивать проект на некоммерческой основе, до тех пор пока они упоминают оригинальное авторство и лицензируют производные работы на аналогичных лицензионных условиях. Все новые работы, основанные на этом проекте, должны иметь эту же лицензию, поэтому все производные работы также должны носить некоммерческий характер.

### Contributing
Мы открыты к предложениям и изменениям, вы можете испрользовать issues или [связаться с нами](https://t.me/semyon_dev).

### Запуск

Необходимо добавить переменные окружения (можно через .env в корне проекта)

Подробнее смотрите пакет config

`go run main.go`

Или скомпилировать в единый бинарник:

`go build main.go`

### Документация к API методам
##### Публичные методы

Протокол: HTTP, формат данных: JSON

<details>
<summary>методы для получения ФСИН учреждений</summary>

Все сразу 
  
```
GET /places
```
Ответ: массив мест:
```
[
     {"_id": "5ed2c5fd0c4a85b90ef09431",
      "name": "ФКУ «ИК № 10 ГУФСИН по Приморскому краю»",
      "type": "Исправительная колония",
      "position": {
        "lat": 43.987453,
        "lng": 132.337293
      },
      "coronavirus": false,
      "number_of_violations": 0},
]
```
Пояснение:
`_id` - уникальный id места (нужен для /places/:id) \
`name` - полное название учреждения \
`type` - тип колонии \
`position` - геолокация \
`coronavirus` - имеется ли информация о коронавирусе \
`number_of_violations` - кол-во нарушений по нашей информации \

Конкретное место:
```
GET /places/<id>
```
Пример ответа для запроса /places/5ed2c5fd0c4a85b90ef09431:
```
{
  "place": {
    "_id": "5ed2c5fd0c4a85b90ef09431",
    "name": "ФКУ «ИК № 10 ГУФСИН по Приморскому краю»",
    "type": "Исправительная колония",
    "position": {
      "lat": 43.987453,
      "lng": 132.337293
    },
    "coronavirus": false,
    "number_of_violations": 0,
    "location": "Михайловский район, пос. Горное",
    "notes": "",
    "phones": [
      "+7 (42346) 3-82-33",
      "+7 (42346) 3-81-31"
    ],
    "hours": "пн-пт 8:00–16:12",
    "website": "http://25.fsin.su/kontaktnaya-informatsiya-po-uchrezhdeniyam-kraya.php?clear_cache=Y",
    "address": "Россия, Приморский край, Михайловский район, поселок Горное, улица Ленина, 25",
    "warning": "",
    "violations": null,
    "corona_violations": null
  }
}
```
Помимо параметров из /places будут:
`location` - местоположение (Город, поселок и тд)
`notes` - заметки учреждения (из википедии)
`phones` - массив телефонов
`hours` - часы работы
`website` - веб сайт
`address` - полный адрес
`warning` - предупреждение (например, место нуждается в проверке)
`violations` - нарушения
`corona_violations` - информация о коронавирусе

</details>

<details>
<summary>методы для получения нарушений</summary>

Нарушения (в том числе по короне) есть только для конкретных учреждений

* получение всех нарушений у которых есть информация по коронавирусу \
`GET /corona_places`

* получение всех нарушений \
`GET /violations`

</details>

<details>
<summary>методы для получения аналитики</summary>
  
* пояснения по разным параметрам (скорее для аналитики) \
`GET /explanations`

* аналитика по разным параметрам (общая статистика) \
`GET /analytics`
Пример ответа:
```
{
  "total_count": 4995,
  "total_count_appeals": 377,
  "total_count_appeals_corona": 105,
  "violations_stats": {
    "communication": {
      "total_count": 1124,
      "total_count_appeals": 967,
      "count_by_years": {
        "2014": 54,
        "2015": 103,
        "2016": 99,
        "2017": 76,
        "2018": 104,
        "2019": 2
      },
      "subcategories": {
        "can_prisoners_submit_complaints": {
          "total_count": 175,
          "total_count_appeals": 366,
          "values": {
            "Да": 84,
            "Затрудняюсь ответить": 107,
            "Нет": 175
          }
        },
       ...
}
```
Параметры:
`total_count` - общее кол-во нарушений по всех заявкам и типам
`total_count_appeals` - общее кол-во заявок
`total_count_appeals_corona` - общее кол-во заявок по коронавирусу
`violations_stats` - категории аналитики, внутри:
`subcategories` - подкатегории

</details>

<details>
<summary>Другое</summary>
  
* получение всех вопросов для создания новых нарушений со стороны клиента `(/form)` \
`GET /formQuestions`

</details>


##### Закрытые методы ([свяжитесь](https://t.me/semyon_dev), чтобы получить доступ)
<details>
<summary>Закрытые методы</summary>
  
* сообщение новых нарушений (форм, заявок)
```
POST /form
place_id string
Параметры нужно получать из GET /formQuestions
```

* создание сообщений по коронавирусу (форм - заявок)
```
POST /form_corona
Параметры:
name_of_fsin string (название МЛС)
place_id string
region string
contacts string
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

</details>
