package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return execute(in, done, stages...)
}

func execute(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	pipe := make(Bi)
	go func() {
		defer close(pipe)
		for {
			select {
			case _, open := <-done:
				if !open {
					return
				}
			case v, open := <-in:
				if !open {
					return
				}
				pipe <- v
			}
		}
	}()

	return execute(stages[0](pipe), done, stages[1:]...)
}
