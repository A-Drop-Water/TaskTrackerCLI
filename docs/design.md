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

* 使用时形如 `.\cli add "Learning Go"`，所以要很好的处理用户的命令行输入。  
    * 错误处理
        * 未收到操作时，打印帮助文档
        * 收到错误操作时，报错
        * 收到操作正确，但参数数量不正确时，报错
    * 正确处理就正常运行
    * 不能使用库是说得到原始的数据全部自己处理是吗，那似乎只需要知道如何获取命令行的字符串即可



### Add tasks
接收用户输入 `.\cli add "Name"`则进行task添加

* 第arg[1]参数得是add
* 参数的长度得是3，即只支持添加 
