# Описание проекта
Go веб-сервер, который реализует функциональность простейшего планировщика задач. Это аналог TODO-листа.

# Список выполенных заданий со звёздочкой. 
- [x] Переменные окружения:
    - TODO_DBFILE
    - TODO_PORT
- [x] Поиск по дате и слову  
- [x] Dokerfile  


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
Таска подготовки docker image - `docker:build`  
Таска запуска docker image    - `docker:run`  

В Taskfile
```yaml
...
vars:
  DOCKER_APP_NAME: scheduler
  DOCKER_IMAGE_NAME: taskscheduler
...
docker:build:
    desc: Build docker image
    cmds:
      - docker build -t {{.DOCKER_IMAGE_NAME}}:v1 .

docker:run:
    desc: Run docker image (default port 7540)
    cmds:
      - docker run --name {{.DOCKER_APP_NAME}} -p 7540:7540 -d {{.DOCKER_IMAGE_NAME}}:v1
      - sleep 2
      - docker ps | grep {{.DOCKER_APP_NAME}}
```

Запуск через go task:
```sh
task docker:build
```
```sh
task docker:run
```