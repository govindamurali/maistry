# maistry
Maistry is a Golang implementation of Worker Pool. Effective when you have limited worker resources but has Jobs/functions to execute across multiple use-cases with customizable requests-per-second.  

## Install

	go get github.com/govindamurali/maistry
## Concepts

* **Job** - Essentially a **func()** to be executed 
* **Worker** - Does a **job**.
* **Job Manager** - This generates the jobs. Takes in two parameters, **jobsPerSecond** and a **name**.
* **Dispatcher** - Dispatches the jobs to the workers. Takes in **maxWorkers** as the input, and a logger interface with standard Error and Info.


## How to use

#### 1. Create dispatcher with workerCount count and start it
```
	dispatcher := maistry.NewDispatcher(workerCount, logger)
	dispatcher.Start() 
```

#### 2. Create Job Managers
```
	jm1:= maistry.NewJobManager(jps1, "manager 1")
	jm2:= maistry.NewJobManager(jps2, "manager 2")
```

#### 3. Run the job managers
```
	jm1.Run()
	jm2.Run()
```

#### 4. Create and push the jobs to the job manager
```
	job1 := maistry.NewJob(
		func() {
			//your function here 
		}, logger)
	jm1.PushJob(job1)

	job2 := maistry.NewJob(func() {
			//your other function here
		}, logger)
	jm2.PushJob(job2)

```

