package main

import (
	"fmt"
	"learning-go/validation"
)

var x = "test1"
var y string = "test2"

type Person struct { // 创建结构体类型
	Name string
	Age  int
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		//TIP <p>To start your debugging session, right-click your code in the editor and select the Debug option.</p> <p>We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.</p>
		fmt.Println("i =", 100/i)
	}

	// 运行基础语法演示
	fmt.Println("\n是否运行基础语法演示？(y/n): ")
	var choice string
	fmt.Scanln(&choice)
	if choice == "y" || choice == "Y" {
		Demo()
	}

	// 运行 Redis 测试
	fmt.Println("\n是否运行 Redis 测试？(y/n): ")
	fmt.Scanln(&choice)
	if choice == "y" || choice == "Y" {
		validation.TestRedisSlotStrategy()
	}

	fmt.Println("\n按回车键退出...")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
}
