# 简单的功能描述和设计思路
## 功能描述
*    Add, Update, and Delete tasks
*    Mark a task as in progress or done
*    List all tasks
*    List all tasks that are done
*    List all tasks that are not done
*    List all tasks that are in progress
## 限制
* Use positional arguments in command line to accept user inputs.
* Use a JSON file to store the tasks in the current directory.
* The JSON file should be created if it does not exist.
* Use the native file system module of your programming language to interact with the JSON file.
* Do not use any external libraries or frameworks to build this project.
* Ensure to handle errors and edge cases gracefully.

## 设计思路
自己的一些话和可能的想法记录一下。
* json文件的定义
  * 发现是不是每次运行程序都要从文件中读取整个的json数据到程序里，然后再对里面的字段进行处理？形如
    * 总共的task数
    * task的一个数组，数组成员是每一个描述task的json字符串信息
  * 即需要的结构体有两个，每次都会从文件中读取整体描述，然后再写入到文件中
    * 整体的描述结构体
    * task结构体本身
* 使用时形如 `.\cli add "Learning Go"`，所以要很好的处理用户的命令行输入。  
    * 错误处理
        * 未收到操作时，打印帮助文档
        * 收到错误操作时，报错
        * 收到操作正确，但参数数量不正确时，报错
    * 正确处理就正常运行
    * 不能使用库是说得到原始的数据全部自己处理是吗，那似乎只需要知道如何获取命令行的字符串即可
* 存储的记录json文件结构
    * 创建时间
        * 类型 time.Time
    * 更新时间
        * 类型 time.Time
    * 完成状态 
        * 用枚举 ?
        * in progress
        * done
        * ToDo 默认
    * 任务名称
    * 任务ID
        * 怎么决定 ? 
        * 比如我完成一个任务后是不是要更新其他的ID? 



### Add tasks
接收用户输入 `.\cli add "Name"`则进行task添加
* 参数规范
    * arg[1]参数得是add
    * ~~ arg[2]参数规范应该以""包围起来,后续可以考虑说没有""情况  ~~
        * 哦这个是命令行输入接收时就已经能够正常处理了
    * 参数的长度得是3，即只支持添加，或者说我们忽略第三个参数之后的系列参数
    * 允许添加同样的任务描述
* 需要用到的技术
    * 接收命令行参数
    * 进行文件的追加写入操作
    * 使用提供的json相关操作实现格式化写入

### List tasks
* 参数规范
  * 可以只有list
  * 或者可以带有状态，只有输入有效状态时才会列出对应的，但是无效状态不会报错
* 实现
  * 就是读取json文件到结构体内，然后直接一一根据指定状态打印即可
### Update tasks
* 参数规范
  * task-cli update 1 "Buy groceries and cook dinner"
  * 感觉这个没什么用啊？ 就是更新一下内容
* 实现
  * 读取json文件到结构体，然后更新结构体的指定内容，即名字，然后再写入到文件
### Delete task
* 参数规范
  * task-cli delete 1
* 实现
  * 和update差不多，只是直接删除对应结构体
  * 注意得遍历更新id和number

### Mark task
* 参数设置
  * task-cli mark-in-progress 1
  * task-cli mark-done 1
* 实现
  * 和前面的更新都差不多，只不过修改的东西不一样
