package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	// 获取命令行参数
	rawArg := os.Args
	if len(rawArg) < 2 {
		// 这里后续可能会答应帮助信息
		fmt.Println("Please enter a command")
		// 直接结束程序
		return
	}

	// 参数合法
	args := rawArg[1:]

	// 进行处理
	switch args[0] {
	case "add":
		// 只传入需要的任务名字参数
		err := AddTask(args)

		if err != nil {
			fmt.Println("error in add:", err)
		} else {
			fmt.Println("Add succeed!")
		}
	case "list":
		// 只需传入list即可
		err := ListAllTasks(args)

		if err != nil {
			fmt.Println("error in list:", err)
		}
	case "update":
		err := UpdateTask(args)

		if err != nil {
			fmt.Println("error in update:", err)
		}
	case "delete":
		err := DeleteTask(args)

		if err != nil {
			fmt.Println("error in delete:", err)
		}
	case "mark-in-progress":
		err := MarkState(args, InProgress)

		if err != nil {
			fmt.Println("error in mark-in-progress:", err)
		}
	case "mark-done":
		err := MarkState(args, Done)

		if err != nil {
			fmt.Println("error in mark-done:", err)
		}
	}

}

// 状态的枚举变量
type State int

// 好像不需要这么进行定义，因为要输出的话也不好搞
const (
	InProgress State = iota
	Done
	Todo
	Unknown
)

// 定义枚举值转成字符串的函数
func (s State) String() string {
	switch s {
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	case Todo:
		return "Todo"

	default:
		return "Unknown"

	}
}

// 文件的地址 定义为常量吧
const JSONFILE = "task.json"

// 任务的结构体
type Task struct {
	// 时间相关
	CreateTime time.Time `json:"create_time"`
	ModifyTime time.Time `json:"modify_time"`
	// 状态
	TaskState State `json:"task_state"`
	// 名称
	TaskName string `json:"task_name"`
	// 任务标号  --->  这个该如何设置 ?
	TaskID int `json:"task_id"`
}

type TaskTracker struct {
	Number int    `json:"number"` // 存储的task总数
	Tasks  []Task `json:"tasks"`  // 存储的每一个task 结构体信息

}

// 从文件中读取整个json文件，返回对应的TaskTracker结构体
func GetTaskTracker() (TaskTracker, error) {
	// 打开json文件
	file, err1 := os.Open(JSONFILE)
	if err1 != nil {
		return TaskTracker{}, err1
	}
	defer file.Close()

	// 读取整个json文件到content
	content, err2 := io.ReadAll(file)

	if err2 != nil {
		return TaskTracker{}, err2
	}
	// 解析对应的json序列到要返回的结构体中 这里就简单考虑，直接返回
	var taskTracker TaskTracker

	// 如果一开始读取的文件为空

	if len(content) == 0 {
		return TaskTracker{}, nil
	}

	err3 := json.Unmarshal(content, &taskTracker)
	if err3 != nil {
		return TaskTracker{}, err3
	}

	return taskTracker, nil

}

// 向文件中写入taskTracker
func WriteTaskTracker(taskTracker TaskTracker) error {
	// 这个权限不清楚重不重要
	file, err := os.OpenFile(JSONFILE, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}
	// 现在写入file
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(taskTracker)

	if err != nil {
		return err
	}

	return nil
}

// 添加任务函数 error 返回值用来
func AddTask(args []string) error {

	// 参数规范
	// 只能有一个名字
	// 仅此而已0
	if len(args) != 2 {
		return errors.New("wrong args in add, only one arg required by add")
	}

	// 参数正确
	// 获取对应的结构体
	taskTracker, err1 := GetTaskTracker()

	if err1 != nil {
		return err1
	}

	// 得到的是空呢 ? 也没关系吧，好像有默认的零值

	// 然后添加相应的数据信息
	// 得创建相应数据信息的结构体

	newTask := Task{
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		//TaskState:  InProgress,
		// 默认为todo
		TaskState: Todo,
		TaskName:  args[1],
		TaskID:    taskTracker.Number + 1,
		//  怎么设置  ? id应为原来存储的数量加1
		// 原来有5个，那我的新的id就是6
	}

	// 更新新的json数据
	// 添加新的task数据 和 number更新
	taskTracker.Tasks = append(taskTracker.Tasks, newTask)
	taskTracker.Number++

	// 接下来文件操作 进行添加
	// 或者说进行重新写入
	// 截断模式打开文件， 不存在就创建  ，读写模式  ,       文件权限设置

	err2 := WriteTaskTracker(taskTracker)

	if err2 != nil {
		return err2
	}

	fmt.Println("NewTask info:", newTask)

	return nil
}

func ListAllTasks(args []string) error {
	// 参数规范就是只有一个参数 list
	// 好的可以指定状态
	if len(args) > 2 {
		return errors.New("wrong args in list")
	}

	// 限定了某一类型
	var s State
	if len(args) == 2 {
		switch args[1] {
		case "done":
			s = Done
		case "in-progress":
			s = InProgress
		case "todo":
			s = Todo
		default:
			s = Unknown
		}
	}

	// 直接列出所有的task
	taskTracker, err := GetTaskTracker()
	if err != nil {
		return err
	}

	if len(taskTracker.Tasks) == 0 {
		fmt.Println("No tasks found!")
	}

	fmt.Println("Tasks:")
	for _, task := range taskTracker.Tasks {

		if len(args) == 1 {
			// 只有一个参数那就都列出来
			fmt.Printf("[%d] |%s| |%s|\n", task.TaskID, task.TaskName, task.TaskState.String())
		} else if task.TaskState == s {
			fmt.Printf("[%d] |%s| |%s|\n", task.TaskID, task.TaskName, task.TaskState.String())
		}
	}

	return nil
}

func UpdateTask(args []string) error {
	// 接收 id 和 新的 string
	// 判断参数合法性 ?
	if len(args) != 3 {
		return errors.New("wrong args in update")
	}

	// 现在知道参数有两个，现在判断一个是不是数字，一个是不是string
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	name := args[2]

	taskTracker, err1 := GetTaskTracker()

	if err1 != nil {
		return err1
	}

	// 看有没有指定id的task
	if id >= 1 && id <= taskTracker.Number {
		// 直接改就行了
		taskTracker.Tasks[id-1].TaskName = name
		// 时间更新
		taskTracker.Tasks[id-1].ModifyTime = time.Now()
	} else {
		return errors.New("task not found")
	}

	// 改完了，写回json

	err2 := WriteTaskTracker(taskTracker)

	if err2 != nil {
		return err2
	}

	return nil

}

func DeleteTask(args []string) error {
	if len(args) != 2 {
		return errors.New("wrong args in delete")
	}

	// 现在知道参数
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	taskTracker, err1 := GetTaskTracker()

	if err1 != nil {
		return err1
	}

	// 看有没有指定id的task
	if id >= 1 && id <= taskTracker.Number {
		// 直接删除对应的task 比如这个task的下标是 id - 1
		// 然后 task [id:]后面全部id减1 number --

		//

		// 边界问题 ? 比如只有一个1       删除后就是0 number变为0   task[1:] 会直接报错吧 ? 特判呗
		// 但好像没有问题
		// 删除的id位置是 id - 1 所以应该是从 id开始减少
		// 这样好像是值传递 ?
		//for _, task := range taskTracker.Tasks[id:] {
		//	task.TaskID--
		//}

		for i := id; i < taskTracker.Number; i++ {
			taskTracker.Tasks[i].TaskID--
		}

		taskTracker.Number--

		taskTracker.Tasks = append(taskTracker.Tasks[:id-1], taskTracker.Tasks[id:]...)

		return WriteTaskTracker(taskTracker)

		//if taskTracker.Number == 1 {
		//	// 只有一个那就是直接删除
		//	// 或者说直接写入一个空的结构体
		//	return WriteTaskTracker(TaskTracker{})
		//} else {
		//	// 没有边界问题
		//	for _, task := range taskTracker.Tasks[id:] {
		//		task.TaskID--
		//	}
		//
		//	// 删除id位置，但是id的位置对应下标是  id - 1
		//	taskTracker.Tasks = append(taskTracker.Tasks[:id-1], taskTracker.Tasks[id:]...)
		//}
	} else {
		return errors.New("task not found")
	}

}

func MarkState(args []string, state State) error {
	if len(args) != 2 {
		return errors.New("wrong args in marking state")
	}

	// 确定参数正确，获取对应状态和id

	id, err := strconv.Atoi(args[1])

	if err != nil {
		return err
	}

	// 好的接下来读取文件
	taskTracker, err1 := GetTaskTracker()

	if err1 != nil {
		return err1
	}

	// 判断id合法性
	if id >= 1 && id <= taskTracker.Number {
		// 合法id 修改state
		taskTracker.Tasks[id-1].TaskState = state
		taskTracker.Tasks[id-1].ModifyTime = time.Now()
	} else {
		return errors.New("task not found")
	}

	// 写回

	err3 := WriteTaskTracker(taskTracker)

	if err3 != nil {
		return err3
	}

	return nil
}
