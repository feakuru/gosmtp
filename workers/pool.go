package workers

import "fmt"

// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func worker(id int, jobs <-chan int, results chan<- int, task func()) {
    for j := range jobs {
        fmt.Println("worker", id, "started  job", j)
        task()
        fmt.Println("worker", id, "finished job", j)
        results <- j * 2
    }
}

func WorkerPool(workersNumber int, task func ()) {

    jobs := make(chan int, workersNumber)
    results := make(chan int, workersNumber)
    defer jobs.close()

    for w := 1; w <= workersNumber; w++ {
        go worker(w, jobs, results, task)
    }

    for j := 1; j <= workersNumber; j++ {
        jobs <- j
    }

    for a := 1; a <= workersNumber; a++ {
        <-results
    }
}
