package main

import "sync"

// сюда писать код
func ExecutePipeline(jobs ...job) {

	wg :=&sync.WaitGroup{}
	in :=make (chan  interface{})
	out :=make (chan  interface{})


	for i:=0; i < len(jobs); i ++{

		in = out
		out = make(chan interface{})

		wg.Add(1)

		go func(i job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			i(in, out)
		}(jobs[i], in, out)
	}
	wg.Wait()
}