package goChanTools
import "testing"

func makeChans(num int) []chan interface{} {
    out := make([]chan interface{}, num)
    for i := range(out) {
        out[i] = make(chan interface{})
    }
    return out
}

func testRead(expected string, expectedOk bool, ch <-chan interface{}, test *testing.T) {
    val, ok := <-ch
    if ok != expectedOk {
        //Didn't get a message and exected one, or got one and didn't expect it.
        test.FailNow()
    }

    sval, convOk := val.(string)
    if convOk != expectedOk {
        //Expected a value, and didn't get one, or got one and didn't expected it.
        test.FailNow()
    }

    //We got what we've expected so far, check that the values line up.

    if expected != sval {
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


    testRead("A", true, c[0], test)
    testRead("B", true, c[0], test)
    testRead("", false, c[0], test) 
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
    testRead("A", true, c[1], test)
    testRead("A", true, c[2], test)

    close(c[0])

    testRead("", false, c[1], test)
    testRead("", false, c[2], test)
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
    testRead("A", true, c[2], test)
    testRead("A", true, c[3], test)

    c[1] <- "B"
    testRead("B", true, c[2], test)
    testRead("B", true, c[3], test)

    close(c[0])

    c[1] <- "C"
    testRead("C", true, c[2], test)
    testRead("C", true, c[3], test)

    close(c[1])
    
    testRead("", false, c[2], test)
    testRead("", false, c[3], test)
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
    
    testRead("A", true, c[3], test)
    testRead("A", true, c[3], test)
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
    
    testRead("A", true, c[1], test)
    testRead("A", true, c[3], test)
    testRead("A", true, c[4], test)
}
