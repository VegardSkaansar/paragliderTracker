package paragliderdb

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ------------------------------------------------------------------------
// -                 -- GLOBAL Datastructures --                          -
//-------------------------------------------------------------------------

//ServiceTime is how long the server has been up
type ServiceTime struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// TrackMeta is the data for a special Igc file
type TrackMeta struct {
	Date        time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	TrackLength float64   `json:"track_length"`
	TrackSrcURL string    `json:"track_src_url"`
}

// TrackID structure for json body when posting
type TrackID struct {
	ID string `json:"id"`
}

// ------------------------------------------------------------------------
// -                 -- GLOBAL Variables --                               -
//-------------------------------------------------------------------------

// StartTime global variable that will be intialised and sets the time from
// when the server starts
var StartTime time.Time

// ------------------------------------------------------------------------
// -                 -- MongoDB --                                        -
//-------------------------------------------------------------------------

// TrackerDB is the interface we use to navigate our database
type TrackerDB interface {
	Init()
	AddURL(url TrackID) error
	GetTrackID() (TrackID, error)
	GetAll()
}

// MongoDB here we store the information about connection
type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

// Init function initilises the mongodb
func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.CollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

}

// AddURL adds a new URL to the database (mongodb)
func (db *MongoDB) AddURL(id TrackID) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(id)

	if err != nil {
		fmt.Printf("Somethings wrong with Insert():%v", err.Error())
		return err
	}

	return nil
}

// GetTrackID gets an id and and we return the trackmeta for this id
func (db *MongoDB) GetTrackID(id string) (TrackID, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	url := TrackID{}
	ok := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"id": id}).One(&url)

	if err != nil {
		ok = false
	}
	return url, ok
}
