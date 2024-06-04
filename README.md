# GoNews

Данный проект является домашней работой на тему создания сервиса-агрегатора новостей. 
Я потратил значительное время на экспериментирование с `docker compose` и пришел к 
структуре сервиса, которую опишу далее. Начну с самого важного - с запуска сервиса.

## Запуск
### Боевой
```sh
docker compose -f compose.prod.yml up 
```
Запуск боевого варианта. Сборка проходит в два этапа - поэтому приложение имеет 
минимальный размер. Однако, кеш на жестком диске хостовой системы не сохранется - 
поэтому повторная сборка потребует снова загрузить все пакеты и собрать все юниты с 
нуля.

### Боевой с тестами
```sh
docker compose -f compose.prod.yml -f compose.prod.add-test.yml up 
```
Запуск боевого варианта с функциональными тестами. О результатах тестов будет 
сообщать юнит `client`. Во время тестирования в боевой среде, не будет просиходить 
точной сверки содержания новостных статей. Но будет проведена проверка формата
отдаваемой информации.

### Разработка
```sh
docker compose -f compose.dev.yml up
```
Запуск варианта для разрабтки. В этом режиме не создаётся отдельный образ, а исходный 
код доставляется через volume. Кроме того, сохраняется кеш и все загруженные
зависимости. Этим вариантом удобно перезапускать сервиса во время разработки.

### Разработка с тестами
```sh
docker compose -f compose.dev.yml -f compose.dev.add-test.yml up 
```
Запуск варианта для разрабтки с функциональными тестами. О результатах тестов 
будет сообщать юнит `client`. Во время тестирования будет развернут минимальный
RSS-сервер `rss.mock`, к которому будет обращаться приложение во время работы - это даёт 
возможность проводить точную сверку ожидаемого списка новостных статей со списком,
полученным, после обработки сервисом.

## Структура проекта
`app` - дирректория с юнитами, непосредственно относящимися к работе сервиса.  
`app/frontend` - исходный код, конфиги и артефакты, относящиеся к работе 
веб-интерфейса приложения.  
`app/backend` - исходный код, конфиги и кеш, относящийся к работе бекенда.  
`app/database` - структура и хранилище, относящиеся к работе базы данных.  
`test` - дирректория с юнитами, относящимися к функциональному тестированию приложения.  
`test/rss.mock` - минимальный сервер, возвращающий RSS-списки новостей при 
обращении к одному из двух URL. Списки новостей заданы в виде XML заранее.  
`test/client` - http-клиент, способный отправить запрос к бекенду и проверить
полученный ответ.

## Структура юнита
Все юниты, перечисленные ранее, имеют подобную структуру:  
`source` - исходный код юнита, который однозначно необходимо версионировать.
Для базы данных это `*.sql`-файл со структурой. Так же здесь находятся `Dockerfile`ы,
если они характерны для юнита.  
`config.*` - настроечные файлы, необходимые для запуска боевого `prod` варианта юнита или 
варианта для разработки `dev`. Конфиги для запуска в режиме `dev` должны весионироваться
для облегчения доставки кода, а конфиги для запуска `prod` - не должны. Однако, для 
облегчения работы с проектом и ввиду того, что проект учебный, дирректория `config.prod`
так же находится под системой контроля версий.  
`volume.*` - файлы, сгенерированные юнитом в боевом режиме `prod` и в режиме разработки
`dev`. Здесь может находиться как хранилище базы данных, так и кеш бекенда/фронтенда.
Файлы в этой дирректори не версионируются.
