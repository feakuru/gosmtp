package workers

import "fmt"
import "time"

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

    // In order to use our pool of workers we need to send
    // them work and collect their results. We make 2
    // channels for this.
    jobs := make(chan int, workersNumber)
    results := make(chan int, workersNumber)

    // This starts up 3 workers, initially blocked
    // because there are no jobs yet.
    for w := 1; w <= workersNumber; w++ {
        go worker(w, jobs, results, task)
    }

    // Here we send 5 `jobs` and then `close` that
    // channel to indicate that's all the work we have.
    for j := 1; j <= workersNumber; j++ {
        jobs <- j
    }
    close(jobs)

    // Finally we collect all the results of the work.
    for a := 1; a <= workersNumber; a++ {
        <-results
    }
}
