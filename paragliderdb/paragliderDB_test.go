package paragliderdb

import (
	"testing"

	"github.com/marni/goigc"
)

func Test_Add_Track(t *testing.T) {

	values := []string{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc", "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}
	testmetaData := igc.ParseLocation(values[0])
	newId := NewUniqueParagliderID()
	GlobalDB.Init()
	GlobalDB.AddURL(testmetaData, newId)

	i := GlobalDB.GetAllID()

	ok := false
	for k := 0; k < len(i); k++ {
		if i[k] == newId {
			ok := true
		}
	}
	if ok == false {
		t.Error("trackmeta didnt get added Internal code error")
	}

}
