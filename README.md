# 简介
goworker是使用GO语言实现的并发worker池，你只需要将并发运行的任务注册到goworker，
goworker会按照既定策略进行并发运行。

# 特色功能
- [并发队列控制](#并发队列控制)
- [任务超时](#任务超时)
- [主动退出](#主动退出)

## <span id="并发队列控制">并发队列控制</span>
100个并发的请求，以10个的并发限制运行
```
func main() {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
	})
	for i := 0; i < 100; i++ {
		j := i
		worker.Add(func() {
			fmt.Println("done ", j)
		})
	}
	worker.IsDone()
}
```

## <span id="任务超时">任务超时</span>
如下代码所示：设置2s超时，超时后会停止未开始的的worker请求
```
func main() {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
		TimeOut:        2000,
	})
	for i := 0; i < 30; i++ {
		j := i
		worker.Add(func() {
			time.Sleep(1 * time.Second)
			fmt.Println("done ", j)
		})
		fmt.Println("add ", i)
	}
	worker.IsDone()
}
```
## <span id="主动退出">主动退出</span>
主动停止，未运行的worker会停止
```
func main() {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
	})
	for i := 0; i < 20; i++ {
		j := i
		worker.Add(func() {
			time.Sleep(1 * time.Second)
			fmt.Println("done ", j)
		})
		fmt.Println("add ", i)
	}
	worker.Exit() // 主动停止
}
```

# next move
- 注册的执行任务支持输入参数
- 部署为后台长期运行的任务，可能要考虑每次提交带ID，以区分不同的并发任务
- 运行的结果返回，整体结果