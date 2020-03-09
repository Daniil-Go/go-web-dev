/*
docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=task_app_list -d mysql

create database task_list_app;

drop database name;

CREATE TABLE `task_list_app`.`lists` (
`id` INT NOT NULL AUTO_INCREMENT,
`name` TEXT NOT NULL,
`description` TEXT NULL,
PRIMARY KEY (`id`),
UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);

CREATE TABLE `task_list_app`.`tasks` (
`id` INT NOT NULL AUTO_INCREMENT,
`list_id` INT NOT NULL,
`text` TEXT NOT NULL,
`complete` BOOLEAN NOT NULL DEFAULT 0,
PRIMARY KEY (`id`),
UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);

drop table tasks;

insert into task_list_app.lists (name, description)
 values ("First List", "First Description"), ("Second List", "Second Description");

insert into task_list_app.tasks (list_id, text, complete)
 values (2, "1 Task", 0), (2, "2 Task", 1), (2, "3 Task", 0), (1, "1 Task", 0), (1, "2 Task", 1);

UPDATE task_list_app.tasks SET complete=1 where id=1;

DELETE FROM task_list_app.tasks WHERE text="1 Task";
 */