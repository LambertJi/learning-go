package main

// 导入包
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ==================== 全局变量声明 ====================
var globalInt int = 100
var globalString string = "全局变量"
var globalBool bool = true

// ==================== 全局常量声明 ====================
const (
	Pi     = 3.14159265358979323846
	MaxRetries = 3
	StatusOK = 200
)

// ==================== 类型声明 ====================
type Number int
type Point struct {
	X, Y float64
}

// ==================== 基础数据类型 ====================
func demonstrateBasicTypes() {
	fmt.Println("\n=== 基础数据类型 ===")

	// 整数类型
	var intVar int = 42
	var int8Var int8 = 127
	var int16Var int16 = 32767
	var int32Var int32 = 2147483647
	var int64Var int64 = 9223372036854775807
	var uintVar uint = 42
	fmt.Printf("整数: int=%d, int8=%d, int16=%d, int32=%d, int64=%d, uint=%d\n",
		intVar, int8Var, int16Var, int32Var, int64Var, uintVar)

	// 浮点类型
	var float32Var float32 = 3.14
	var float64Var float64 = 3.141592653589793
	fmt.Printf("浮点: float32=%f, float64=%f\n", float32Var, float64Var)

	// 字符串
	var stringVar string = "Hello, Go!"
	var runeVar rune = 'A' // Unicode 字符
	var byteVar byte = 'B'  // ASCII 字符
	fmt.Printf("字符串: %s, rune=%c, byte=%c\n", stringVar, runeVar, byteVar)

	// 布尔类型
	var boolVar bool = true
	fmt.Printf("布尔: %t\n", boolVar)

	// 常量
	fmt.Printf("常量: Pi=%.2f, MaxRetries=%d, StatusOK=%d\n", Pi, MaxRetries, StatusOK)
}

// ==================== 变量声明方式 ====================
func demonstrateVariableDeclaration() {
	fmt.Println("\n=== 变量声明方式 ===")

	// 完整声明
	var name string = "jiwei"
	var age int = 25
	fmt.Printf("完整声明: name=%s, age=%d\n", name, age)

	// 类型推断
	var city = "Beijing"
	var score = 95.5
	fmt.Printf("类型推断: city=%s, score=%.1f\n", city, score)

	// 短变量声明（最常用）
	country := "China"
	count := 100
	fmt.Printf("短变量声明: country=%s, count=%d\n", country, count)

	// 多变量声明
	x, y := 10, 20
	fmt.Printf("多变量: x=%d, y=%d\n", x, y)

	// 因果赋值（用于交换）
	x, y = y, x
	fmt.Printf("交换后: x=%d, y=%d\n", x, y)

	// 零值
	var zeroInt int
	var zeroString string
	var zeroBool bool
	var zeroFloat float64
	fmt.Printf("零值: int=%d, string='%s', bool=%t, float=%f\n",
		zeroInt, zeroString, zeroBool, zeroFloat)
}

// ==================== 字符串处理 ====================
func demonstrateStringHandling() {
	fmt.Println("\n=== 字符串处理 ===")

	str := "Hello, Go!"

	// 字符串长度
	fmt.Printf("长度: %d (字节), %d (字符数)\n", len(str), len([]rune(str)))

	// 字符串拼接
	str1 := "Hello"
	str2 := "World"
	result := str1 + ", " + str2 + "!"
	fmt.Printf("拼接: %s\n", result)

	// 字符串格式化
	name := "jiwei"
	age := 25
	formatted := fmt.Sprintf("姓名: %s, 年龄: %d", name, age)
	fmt.Printf("格式化: %s\n", formatted)

	// 字符串切片（子串）
	fullString := "Hello, World!"
	substr := fullString[0:5] // "Hello"
	fmt.Printf("子串: %s\n", substr)

	// 字符串遍历（字节）
	fmt.Print("按字节遍历: ")
	for i := 0; i < len(fullString); i++ {
		fmt.Printf("%c ", fullString[i])
	}
	fmt.Println()

	// 字符串遍历（字符）
	fmt.Print("按字符遍历: ")
	for _, char := range fullString {
		fmt.Printf("%c ", char)
	}
	fmt.Println()

	// 字符串转换
	numStr := "42"
	num := 0
	fmt.Sscanf(numStr, "%d", &num)
	fmt.Printf("字符串转数字: %s -> %d\n", numStr, num)
}

// ==================== 数组和切片 ====================
func demonstrateArraysAndSlices() {
	fmt.Println("\n=== 数组和切片 ===")

	// 数组（固定长度）
	var arr [5]int = [5]int{1, 2, 3, 4, 5}
	fmt.Printf("数组: %v, 长度: %d\n", arr, len(arr))

	// 切片（动态长度）
	var slice []int = []int{1, 2, 3}
	fmt.Printf("切片: %v, 长度: %d, 容量: %d\n", slice, len(slice), cap(slice))

	// 使用 make 创建切片
	slice2 := make([]int, 3, 5) // 长度3，容量5
	fmt.Printf("make创建的切片: %v, 长度: %d, 容量: %d\n", slice2, len(slice2), cap(slice2))

	// 切片操作
	slice = append(slice, 4, 5) // 追加元素
	fmt.Printf("追加后: %v\n", slice)

	slice = append(slice, slice2...) // 追加切片
	fmt.Printf("追加切片后: %v\n", slice)

	// 切片截取
	subSlice := slice[1:4]
	fmt.Printf("截取 slice[1:4]: %v\n", subSlice)

	// 切片复制
	newSlice := make([]int, len(slice))
	copy(newSlice, slice)
	fmt.Printf("复制的切片: %v\n", newSlice)

	// 删除元素（通过append）
	index := 2
	slice = append(slice[:index], slice[index+1:]...)
	fmt.Printf("删除索引%d后的切片: %v\n", index, slice)
}

// ==================== Map ====================
func demonstrateMaps() {
	fmt.Println("\n=== Map ===")

	// 创建 map
	person := make(map[string]int)
	person["Alice"] = 25
	person["Bob"] = 30
	fmt.Printf("Map: %v\n", person)

	// 字面量创建
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
		"Charlie": 35,
	}
	fmt.Printf("字面量创建: %v\n", ages)

	// 读取值
	age := ages["Alice"]
	fmt.Printf("Alice的年龄: %d\n", age)

	// 检查键是否存在
	value, exists := ages["David"]
	fmt.Printf("David存在? %t, 值: %d\n", exists, value)

	// 删除键
	delete(ages, "Bob")
	fmt.Printf("删除Bob后: %v\n", ages)

	// 遍历 map
	fmt.Print("遍历map: ")
	for key, value := range ages {
		fmt.Printf("%s:%d ", key, value)
	}
	fmt.Println()
}

// ==================== 条件判断 ====================
func demonstrateConditionals() {
	fmt.Println("\n=== 条件判断 ===")

	// if 语句
	score := 85
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// if 带初始化语句
	if num := 42; num%2 == 0 {
		fmt.Printf("%d 是偶数\n", num)
	}

	// switch 语句
	day := time.Now().Weekday()
	switch day {
	case time.Monday:
		fmt.Println("今天是星期一")
	case time.Tuesday:
		fmt.Println("今天是星期二")
	case time.Wednesday:
		fmt.Println("今天是星期三")
	default:
		fmt.Printf("今天是其他星期 (%s)\n", day)
	}

	// switch 无条件（替代 if-else）
	hour := time.Now().Hour()
	switch {
	case hour < 12:
		fmt.Println("早上好")
	case hour < 18:
		fmt.Println("下午好")
	default:
		fmt.Println("晚上好")
	}
}

// ==================== 循环 ====================
func demonstrateLoops() {
	fmt.Println("\n=== 循环 ===")

	// for 循环（标准形式）
	fmt.Print("标准for循环: ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// while 风格（Go 只有 for）
	fmt.Print("while风格: ")
	i := 0
	for i < 5 {
		fmt.Printf("%d ", i)
		i++
	}
	fmt.Println()

	// 无限循环
	fmt.Print("无限循环(前3次): ")
	count := 0
	for {
		fmt.Printf("%d ", count)
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Println()

	// 遍历切片
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Print("遍历切片: ")
	for index, value := range numbers {
		fmt.Printf("[%d]=%d ", index, value)
	}
	fmt.Println()

	// 遍历 map
	capitals := map[string]string{"China": "Beijing", "Japan": "Tokyo", "USA": "Washington"}
	fmt.Print("遍历map: ")
	for country, capital := range capitals {
		fmt.Printf("%s:%s ", country, capital)
	}
	fmt.Println()

	// 跳过和继续
	fmt.Print("continue示例(只输出偶数): ")
	for i = 0; i < 10; i++ {
		if i%2 != 0 {
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()
}

// ==================== 函数 ====================
func demonstrateFunctions() {
	fmt.Println("\n=== 函数 ===")

	// 基本函数调用
	result := add(10, 20)
	fmt.Printf("add(10, 20) = %d\n", result)

	// 多返回值
	quotient, remainder := divide(17, 5)
	fmt.Printf("divide(17, 5) = %d 余 %d\n", quotient, remainder)

	// 可变参数
	sum := sum(1, 2, 3, 4, 5)
	fmt.Printf("sum(1,2,3,4,5) = %d\n", sum)

	// 闭包
	adder := makeAdder(10)
	fmt.Printf("闭包 makeAdder(10)(5) = %d\n", adder(5))

	// 函数作为参数
	numbers := []int{1, 2, 3, 4, 5}
	result = applyOperation(numbers, func(x, y int) int { return x + y })
	fmt.Printf("应用操作: %d\n", result)
}

// 基本函数
func add(a, b int) int {
	return a + b
}

// 多返回值函数
func divide(a, b int) (int, int) {
	quotient := a / b
	remainder := a % b
	return quotient, remainder
}

// 可变参数函数
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// 闭包函数
func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

// 高阶函数
func applyOperation(numbers []int, op func(int, int) int) int {
	if len(numbers) == 0 {
		return 0
	}
	result := numbers[0]
	for _, num := range numbers[1:] {
		result = op(result, num)
	}
	return result
}

// ==================== 方法和接口 ====================
// 定义结构体
type Rectangle struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

// 定义接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 的方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 的方法
func (c Circle) Area() float64 {
	return Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * Pi * c.Radius
}

// 使用接口的函数
func printShapeInfo(s Shape) {
	fmt.Printf("面积: %.2f, 周长: %.2f\n", s.Area(), s.Perimeter())
}

func demonstrateMethodsAndInterfaces() {
	fmt.Println("\n=== 方法与接口 ===")

	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 4}

	// 调用方法
	fmt.Printf("矩形: 宽=%.2f, 高=%.2f\n", rect.Width, rect.Height)
	fmt.Printf("圆形: 半径=%.2f\n", circle.Radius)

	// 使用接口
	fmt.Print("矩形信息: ")
	printShapeInfo(rect)

	fmt.Print("圆形信息: ")
	printShapeInfo(circle)

	// 类型断言
	var shape Shape = rect
	if r, ok := shape.(Rectangle); ok {
		fmt.Printf("类型断言成功: 这是一个矩形, 宽=%.2f, 高=%.2f\n", r.Width, r.Height)
	}

	// 空接口
	var anything interface{} = 42
	fmt.Printf("空接口值: %v (类型: %T)\n", anything, anything)
}

// ==================== 错误处理 ====================
func demonstrateErrorHandling() {
	fmt.Println("\n=== 错误处理 ===")

	// 基本错误处理
	result, err := divideWithError(10, 2)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %d\n", result)
	}

	// 错误处理
	result, err = divideWithError(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %d\n", result)
	}

	// 自定义错误
	err = validateAge(15)
	if err != nil {
		fmt.Printf("验证错误: %v\n", err)
	}

	// panic 和 recover
	demonstratePanicRecover()
}

func divideWithError(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("不能除以零")
	}
	return a / b, nil
}

func validateAge(age int) error {
	if age < 18 {
		return fmt.Errorf("年龄 %d 小于18岁，不允许访问", age)
	}
	return nil
}

func demonstratePanicRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到panic: %v\n", r)
		}
	}()

	fmt.Println("在panic之前")
	// 这里会触发 panic，但会被 recover 捕获
	panic("这是一个故意的panic")
	fmt.Println("这行不会执行")
}

// ==================== defer 和 init ====================
var initOrder int

func init() {
	initOrder++
	fmt.Printf("Init函数 #%d 执行 (包级别)\n", initOrder)
}

func demonstrateDefer() {
	fmt.Println("\n=== Defer ===")

	// defer 的执行顺序（LIFO）
	defer fmt.Println("Defer #3: 最后执行")
	defer fmt.Println("Defer #2: 第二执行")
	defer fmt.Println("Defer #1: 首先执行")

	// defer 用于资源清理
	file := mockOpenFile("example.txt")
	defer mockCloseFile(file)

	fmt.Println("文件操作中...")

	// defer 用于捕获返回值
	result := deferReturnValue()
	fmt.Printf("defer返回值修改: %d\n", result)
}

func mockOpenFile(filename string) string {
	fmt.Printf("打开文件: %s\n", filename)
	return filename
}

func mockCloseFile(filename string) {
	fmt.Printf("关闭文件: %s\n", filename)
}

func deferReturnValue() (result int) {
	defer func() {
		result *= 2 // defer 可以修改返回值
	}()
	result = 10
	return result
}

// ==================== 并发与协程 ====================
func demonstrateConcurrency() {
	fmt.Println("\n=== 并发与协程 ===")

	// 基本 goroutine
	fmt.Println("启动goroutine...")
	go sayHello("World")
	go sayHello("Go")
	time.Sleep(100 * time.Millisecond) // 等待goroutine完成

	// Channel 基本使用
	fmt.Println("\nChannel示例:")
	ch := make(chan string)
	go sendMessage(ch, "Hello from channel!")
	msg := <-ch
	fmt.Printf("收到消息: %s\n", msg)

	// Buffered channel
	fmt.Println("\nBuffered channel示例:")
	bufferedCh := make(chan int, 2)
	bufferedCh <- 1
	bufferedCh <- 2
	fmt.Printf("从buffered channel读取: %d, %d\n", <-bufferedCh, <-bufferedCh)

	// 方向 selective
	fmt.Println("\nChannel方向:")
	ch2 := make(chan int, 2)
	ch2 <- 10
	ch2 <- 20
	receiveOnly(ch2)

	// select
	fmt.Println("\nSelect示例:")
	ch3 := make(chan string)
	ch4 := make(chan string)
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch3 <- "消息1"
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch4 <- "消息2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch3:
			fmt.Printf("从ch3收到: %s\n", msg)
		case msg := <-ch4:
			fmt.Printf("从ch4收到: %s\n", msg)
		case <-time.After(200 * time.Millisecond):
			fmt.Println("超时")
		}
	}

	// WaitGroup
	fmt.Println("\nWaitGroup示例:")
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	fmt.Println("所有worker完成")

	// Mutex
	fmt.Println("\nMutex示例:")
	demonstrateMutex()
}

func sayHello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

func sendMessage(ch chan<- string, msg string) {
	ch <- msg
}

func receiveOnly(ch <-chan int) {
	val := <-ch
	fmt.Printf("从只读channel读取: %d\n", val)
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d 开始工作\n", id)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Worker %d 完成工作\n", id)
}

func demonstrateMutex() {
	var counter int
	var mutex sync.Mutex

	for i := 1; i <= 5; i++ {
		go func(id int) {
			mutex.Lock()
			counter++
			fmt.Printf("Goroutine %d: counter = %d\n", id, counter)
			mutex.Unlock()
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("最终counter值: %d\n", counter)
}

// ==================== 主函数 ====================
func Demo() {
	fmt.Println("========== Go 基础语法演示 ==========")

	demonstrateBasicTypes()
	demonstrateVariableDeclaration()
	demonstrateStringHandling()
	demonstrateArraysAndSlices()
	demonstrateMaps()
	demonstrateConditionals()
	demonstrateLoops()
	demonstrateFunctions()
	demonstrateMethodsAndInterfaces()
	demonstrateErrorHandling()
	demonstrateDefer()
	demonstrateConcurrency()

	fmt.Println("\n========== 演示完成 ==========")
}
