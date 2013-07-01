package goChanTools
import "testing"

func makeChans(num int) []chan string {
    out := make([]chan string, num)
    for i := range(out) {
        out[i] = make(chan string)
    }
    return out
}

func testRead(expected string, ch <-chan string, test *testing.T) {
    if expected != <-ch {
        test.FailNow()
    }
}

func TestManyToOne(test *testing.T) {
    //
    // c1
    //   >---c0
    // c2
    //
    c := makeChans(3)

    mto := NewManyToOneChan(c[0])
    mto.AddInputChan(c[1])
    mto.AddInputChan(c[2])

    go func () {
        c[1] <- "A"
        close(c[1])

        c[2] <- "B"
        close(c[2]) 
    } ()

    val, ok := <- c[0]
    if val != "A" || ok != true {
        test.FailNow()
    }

    val, ok = <- c[0]
    if val != "B" || ok != true {
        test.FailNow()
    }

    val, ok = <- c[0]
    if val != "" || ok != false {
        test.FailNow()
    }
}

func TestOneToMany(test *testing.T) {
    //
    //         c1
    // c0 ---<
    //         c2
    //
    c := makeChans(3)

    otm := NewOneToManyChan(c[0])
    otm.AddOutputChan(c[1])
    otm.AddOutputChan(c[2])
    
    c[0] <- "A"
    testRead("A", c[1], test)
    testRead("A", c[2], test)

    close(c[0])

    val, ok := <- c[1]
    if val != "" || ok != false {
        test.FailNow()
    }

    val, ok = <- c[2]
    if val != "" || ok != false {
        test.FailNow()
    }
}

func TestManyToMany(test *testing.T) {
    //
    // c0    c2
    //   >--<
    // c1    c3
    //
    c := makeChans(4)
 
    mtm := NewManyToManyChan()
    mtm.AddInputChan(c[0])
    mtm.AddInputChan(c[1])
    mtm.AddOutputChan(c[2])
    mtm.AddOutputChan(c[3])

    c[0] <- "A"
    testRead("A", c[2], test)
    testRead("A", c[3], test)

    c[1] <- "B"
    testRead("B", c[2], test)
    testRead("B", c[3], test)

    close(c[0])

    c[1] <- "C"
    testRead("C", c[2], test)
    testRead("C", c[3], test)

    close(c[1])
    
    val, ok := <- c[2]
    if val != "" || ok != false {
        test.FailNow()
    }
    
    val, ok = <-c[3]
    if val != "" || ok != false {
        test.FailNow()
    }
}

func TestOneToManyToOneDuplication(test *testing.T) {
    //
    //      c1
    // c0 -<  >- c3
    //      c2
    //
    c := makeChans(4)

    otm := NewOneToManyChan(c[0])
    otm.AddOutputChan(c[1])
    otm.AddOutputChan(c[2])

    mto := NewManyToOneChan(c[3])
    mto.AddInputChan(c[1])
    mto.AddInputChan(c[2])

    c[0] <- "A"
    
    testRead("A", c[3], test)
    testRead("A", c[3], test)
}

func TestOneToManyTree(test *testing.T) {
    //
    //      c1
    // c0 -<     c3
    //      c2 -<
    //           c4
    //
    c := makeChans(5)

    otm1 := NewOneToManyChan(c[0])
    otm1.AddOutputChan(c[1])
    otm1.AddOutputChan(c[2])

    otm2 := NewOneToManyChan(c[2])
    otm2.AddOutputChan(c[3])
    otm2.AddOutputChan(c[4])

    c[0] <- "A"
    
    testRead("A", c[1], test)
    testRead("A", c[3], test)
    testRead("A", c[4], test)
}
