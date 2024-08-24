# Описание проекта
Go веб-сервер, который реализует функциональность простейшего планировщика задач. Это аналог TODO-листа.

# Список выполенных заданий со звёздочкой. 
Если их нет, напишите, что задания повышенной трудности не выполнялись;

# Инструкция по запуску кода локально
: дополнительные флаги, примеры .env и так далее. Напишите, какой адрес следует указывать в браузере;
Переменные окружения:  
`TODO_DBFILE` - путь к базе данных (по умолчанию `scheduler.db`)
`TODO_PORT` - порт вебсервера (по умолчанию 7540)

# Инструкция по запуску тестов.
Для удобства запуск тестов реализованы через Taskfile.  
Установка go-task тут https://taskfile.dev/installation/

Посмотреть все доступные таски:  
`task --list-all`

Запуск тестов - `task test-all`


Укажите, какие параметры в tests/settings.go следует использовать;

# Инструкция по сборке и запуску проекта через докер 
если вы это сделали;
