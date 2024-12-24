package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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

	// fmt.Println("arg size is ", len(args))

	// 进行处理
	switch args[0] {
	case "add":
		// 只传入需要的任务名字参数
		err := AddTask(args)

		if err != nil {
			fmt.Println("error :", err)
		} else {
			fmt.Println("Add succeed!")
		}
	case "list":
		// 只需传入list即可
		err := ListAllTasks(args)

		if err != nil {
			fmt.Println("error :", err)
		}

	}

	//fmt.Println(args)
}

// 状态的枚举变量
type State int

// 好像不需要这么进行定义，因为要输出的话也不好搞
const (
	InProgress State = iota
	Done
)

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
		TaskState:  InProgress,
		TaskName:   args[0],
		TaskID:     taskTracker.Number + 1,
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
	file, err := os.OpenFile("task.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)

	if err != nil {
		// fmt.Println("Error when OpenFile : ", err)
		// return errors.New("add task open file fail")
		return fmt.Errorf("failed to open file 'task.json' for appending: %w", err)
	}

	defer file.Close()

	// 写入文件
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err2 := encoder.Encode(taskTracker)

	if err2 != nil {
		return err2
	}

	fmt.Println("NewTask info:", newTask)

	return nil
}

func ListAllTasks(args []string) error {
	// 参数规范就是只有一个参数 list
	if len(args) != 1 {
		return errors.New("too many args in list, only one arg required by list")
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
		fmt.Printf("[%d]: %s\n", task.TaskID, task.TaskName)
	}

	return nil
}
