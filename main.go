package main

// 因为体量很小,所以所有的配置文件都是写死的捏
func main() {
	app := WireApp()
	//启动app
	app.Run()
}
