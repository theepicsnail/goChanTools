package goChanTools
import "sync"

type ManyToOneChan struct {
    sync.WaitGroup
    sync.Once
    outChan chan<- string
}

func NewManyToOneChan(outputChan chan<- string) *ManyToOneChan {
    mto := new(ManyToOneChan)
    mto.outChan = outputChan
    return mto
}

func (m *ManyToOneChan) AddInputChan(ch <-chan string) {
    go func() {
        m.Add(1)
        for line := range(ch) {
            m.outChan <- line
        }
        m.Done()
        m.Do(func() {
            m.Wait()
            close(m.outChan)
        })
    }()
}

type OneToManyChan struct {
    inChan <-chan string
    outChans [] chan<- string
}

func NewOneToManyChan(srcChan <-chan string) *OneToManyChan {
    o := new(OneToManyChan)
    o.inChan = srcChan
    o.start()
    return o
}

func (o *OneToManyChan) start() {
     go func() {
        for line := range(o.inChan) {
            for _, ch := range(o.outChans) {
                go func(line string, ch chan<- string){
                    ch <- line
                }(line, ch)
            }
        }
        for _, ch := range(o.outChans) {
            close(ch)
        }
    }()
}

func (o *OneToManyChan) AddOutputChan(ch chan<- string) {
    o.outChans = append(o.outChans, ch)
}

type ManyToManyChan struct {
    ManyToOneChan
    OneToManyChan
}
func NewManyToManyChan () *ManyToManyChan {
    middleChan := make(chan string)
    mtm := new(ManyToManyChan)
    mtm.outChan = middleChan
    mtm.inChan = middleChan
    mtm.start()
    return mtm
}

