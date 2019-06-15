package database

import "fmt"
import "testing"
import "time"
import "encoding/json"

import "github.com/fatih/structs"

type Job struct{
	Name string
	ID string
	Updated time.Time
}

func TestStructMapMarshalling(t *testing.T) {
	j := Job{
		Name: "foo",
		ID: "1",
		Updated: time.Now(),
	}
	var m map[string]interface{}
	m = structs.Map(j)
	fmt.Printf(m["Name"].(string))

	ba,_ := json.Marshal(m)
	str := string(ba)
	fmt.Printf(str)
}

func TestSourcerSinkerPattern(t *testing.T) {
	// reader
	ch := source()

	// writer
	sink(ch)

}
func source() chan string{
	ch := make(chan string)
	go func(){
		for i := 0; i < 100; i++{
			ch <- "test"
		}
		close(ch)
	}()
	return ch
}
func sink(source chan string){
	for x := range source{
		fmt.Println(x)
	}
}
