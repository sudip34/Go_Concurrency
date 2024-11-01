GO Concurrency
===

Share memory by Communicating

Don't over-engineer things by using shared memory and complicated, error-prone synchronization primitivies, use message-passing between GoRoutines so variables and data can be used in the appropriate sequence.
***GOLDEN RULE*** for ***WHEN TO USE***

> if you don't need it, don't use it.
> Keep your application's complesity to an absolute minimum; it's easier to write, easier to understand, and easier to maintain


## Requiered
 - GO
 - Make
 - IDE(VS_CODE) + GO (extension) 

 ### VS-CODE configuration
 - insatll GO
 - `crtl` + `shift` +`p` 
 - write `go:insatll/ Update tools`
 - extensions : `GO`, `gotemplate-syntax`

 ### Make
 - install Make


 ## comand to build a project in go

  `go mod init <project-name>`

- to update a project/ bring all dependencies
  `go mod tidy`  

- to run a go app
> go run . 

or
>go run <main.go>

- for test:
> go test .
 ### GoRutines

 - Running things in the background
 - `main` function itself a GoRoutine
 - GoRoutine is very light weight Thread. They are all managaed (with schedualer) and they need very very low memory. Where in Java Platform Thread needs ~ 2MB
 - GoRoutine is not platform Thread.
 - GoRoutine can not communicate with eatch other like channels

 ### WaitGropus instaed of time.sleep(*)
 - we need to make the main thread to wait for the  another goRoutine to get the output
 - doing this with `time.Sleep(1 * time.Second)`  is the worst way to do it
 - need to do with `WaitGropus`

 how to use WaitGroup:

 > the function will be called in a GoRoutine it should take a pointer of `sync.WaitGroup`

 > `sync.WaitGroup` need to use `WaitGroup.Add(No of elements in slice)` after the `slice`  

 > then the function will be called inside a `GoRutine`

 > And after the `goRutine` we need to use `waitGroup.wait` so that `main()` waits for the `GoRoutine` to execute.
 
 > we should use a `pointer` of `sync.WaitGroup` as according to Go-documentation we should not copy a `sync.WaitGroup`

 > We need to use `defer` + `waitGroup.Done()` in begining of the `function` that we are going to call inside `GoRoutine`

 >`defer` + `waitGroup.Done()`  deduct from the Value added in `.Add()`  and at `wait()`  the value of `No of elements in slice` in `.Add(0)` is 0 again.

 > So, After the `wait()` if we are going the use the  above same `function`, which takes a `WaitGroup`, we need to use `WaitGroup.Add(number)`

 ```
 func printSomething(s string, waitGroup *sync.WaitGroup) {
    //use defer.Done() so that this function should  finish execution first and only then anything else will be executed
    defer waitGroup.Done()
	fmt.Println(s)
}
 func main() {
 	var waitGroup sync.WaitGroup

	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"pi",
		"zeta",
	}
	// so as the slice has 6 elements so waitGroup need to wait for 6
	waitGroup.Add(len(words))

	for index, x := range words {
		go printSomething(fmt.Sprintf("%d: %s", index, x), &waitGroup)
	}

	waitGroup.Wait() // using instead of:    time.Sleep(1 * time.Second) 

    waitGroup.Add(1) // need to add 1 for the next printSomething(string, waitGroup)
    printSomething("this is the second thing to be printed", &waitGroup)

 }
 ```

## RACE conditions, MUTEXES, and CHANNELS

***to check race condition***

`go run -race .`

`go test -race .`

- install `gcc` compiler to check `-race`
- run in cmd: `set CGO_ENABLED=1`

***MUTEX*** - mutual exclusion
- relatively simple to use
- Dealing with sheared resources and concurrent/parallel goroutines 
- Lock/Unlock

> TO get rid of race condition arisen with go GoRoutine can be solved with sync.Mutex along with sync.WaitGroup

> after defer sync.WaitGroup we lock the function and unlock it

```
defer waitGroup.Done()

mutex.Lock()
something that will performed 
mutex.Unlock()
```
example of `sync.Mutex`
```
package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string, m *sync.Mutex) {
	defer wg.Done()
	m.Lock()
	msg = s
	m.Unlock()
}

func main() {
	msg = "Hello, world!"

	var mutex sync.Mutex

	wg.Add(2)

	go updateMessage("Hello, universe!", &mutex)

	go updateMessage("Hello, cosmos!", &mutex)
	wg.Wait()

	fmt.Print(msg)

}
```

### Channels 
- Channels are a means of having GoRoutines share data
- They can talk to each other
- Channels can be Unidirectional and bidirectional

> GOLDEN RULE:  We must close a chan after the use

> `chan chan` is a mean to close a `chan` 

> we use `select` statement only for channels

```
select {
		//we sent something to the data channel
case pizzaMaker.data <- *currentPizza:

case quitChan := <-pizzaMaker.quit:

}
```
