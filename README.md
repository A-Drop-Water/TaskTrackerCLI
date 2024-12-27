# TaskTrackerCLI
First small Go project.
From https://roadmap.sh/projects/task-tracker

## Use for the this simple project

* build the project
  * `go build`
* add tasks 
  * `./TaskTracker add "Your task name or desc"`
* list tasks
  * `./TaskTracker list` or `./TaskTracker list [done|in-progress|todo]`
* delete task
  * `./TaskTracker delete id(The id of the task)`
* update task name
  * `./TaskTracker update id "New task name or desc"`
* change state
  * `./TaskTracker mark-done id`
  * `./TaskTracker mark-in-progress id`
