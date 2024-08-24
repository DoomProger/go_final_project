# Описание проекта
Go веб-сервер, который реализует функциональность простейшего планировщика задач. Это аналог TODO-листа.

# Список выполенных заданий со звёздочкой. 
Переменные окружени:
- TODO_DBFILE
- TODO_PORT

Dokerfile


# Инструкция по запуску кода локально
Переменные окружения:  
`TODO_DBFILE` - путь к базе данных (по умолчанию `scheduler.db`)  
`TODO_PORT` - порт вебсервера (по умолчанию 7540)  

# Инструкция по запуску тестов.
Для удобства запуск тестов реализованы через Taskfile.  
Установка go-task тут https://taskfile.dev/installation/

Посмотреть все доступные таски:  
```sh
task --list-all
```

Запуск тестов: 
```sh
task test-all
```


Укажите, какие параметры в tests/settings.go следует использовать;

# Инструкция по сборке и запуску проекта через докер 
Таска подготовки docker image - `build-docker`  
Таска запуска docker image    - `run-docker`  

```yaml
build-docker:
    desc: Build docker image
    cmds:
      - docker build -t taskscheduler:v1 .

run-docker:
    desc: Run docker image on port 7540
    cmds:
      - docker run taskscheduler:v1 -p 7540:7540
```

Запуск через go task:
```sh
task build-docker
```
```sh
task run-docker
```