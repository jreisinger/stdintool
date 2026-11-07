package stdintool

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
)

type Parser interface {
	Parse(line string) Task
}

type Task interface {
	Process()
	Print()
}

func Run(p Parser, nWorkers int) {
	var wg sync.WaitGroup
	in := make(chan Task)
	out := make(chan Task)

	wg.Go(func() {
		defer close(in)
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			in <- p.Parse(s.Text())
		}
		if err := s.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading from stdin: %v\n", err)
			os.Exit(1)
		}
	})

	if nWorkers < 1 {
		nWorkers = runtime.NumCPU()
	}

	for range nWorkers {
		wg.Go(func() {
			for t := range in {
				t.Process()
				out <- t
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for t := range out {
		t.Print()
	}
}
