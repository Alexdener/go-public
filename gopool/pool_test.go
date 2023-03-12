package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	p := NewPool(5)
	for i := 0; i < 20; i++ {
		ii := i
		p.NewTask(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(time.Second)
			}
			return nil
		})
	}
	p.Wait()
}
