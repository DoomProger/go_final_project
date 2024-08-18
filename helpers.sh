export TODO_DBFILE=scheduler.db
export TODO_PORT=8080


# test 3
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests

rm -f /home/nkorolev/repos/yandexpractikum/go_final_project/scheduler.db
rm -f /home/nkorolev/repos/yandexpractikum/scheduler.db
