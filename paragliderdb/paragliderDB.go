package paragliderdb

import (
	"fmt"
	"log"
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
// -                 -- GLOBAL Const Variables --                         -
//-------------------------------------------------------------------------

// MAXLENGTHID variable for the lenght of the random crypto generator for
// paragliderID, this will make it easier to track the length
const MAXLENGTHID = 6

// ------------------------------------------------------------------------
// -                 -- GLOBAL Variables --                               -
//-------------------------------------------------------------------------

// StartTime global variable that will be intialised and sets the time from
// when the server starts
var StartTime time.Time

// GlobalDB is the global variable that will store and find all our
// information from the clients requests
var GlobalDB TrackerDB

// ------------------------------------------------------------------------
// -                 -- Interfaces --                                     -
//-------------------------------------------------------------------------

// TrackerDB is the interface we use to navigate our database
// this is the interface that handles all our request
// and returns what we need or store what we asked for
type TrackerDB interface {
	Init()
	AddURL(meta TrackMeta, newID TrackID) error
	GetTrackID(id string) (TrackID, bool)
	GetAllID() []TrackID
	CheckIfURLIsAlreadyTracked(url string) bool
}

// ------------------------------------------------------------------------
// -                 -- MongoDB --                                        -
//-------------------------------------------------------------------------

// MongoDB here we store the information about connection
type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

// this is just used for making the collection structure
type collection struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ParagliderID TrackID
	Track        TrackMeta
}

// Init function initilises the mongodb
func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

}

// AddURL adds a new URL to the database (mongodb)
func (db *MongoDB) AddURL(meta TrackMeta, newID TrackID) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	col := collection{
		bson.NewObjectId(),
		newID,
		meta,
	}

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(col)

	if err != nil {
		fmt.Printf("Somethings wrong with Insert():%v", err.Error())
		return err
	}

	return nil
}

// GetTrackID gets an id and and we return if this for this id
func (db *MongoDB) GetTrackID(id string) (TrackID, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	idT := TrackID{}
	ok := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(nil).Select(bson.M{"paragliderid": bson.M{"$elemMatch": bson.M{"id": id}}}).One(&idT)

	if err != nil {
		ok = false
	}
	return idT, ok
}

// CheckIfURLIsAlreadyTracked handles the situation when mulitple
// people post the same track
func (db *MongoDB) CheckIfURLIsAlreadyTracked(url string) bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	var idT []TrackID

	// we go to the db and the collection and check if we already have this tracksrcurl
	// if not we wil get a empty body and this equals 0 when we ask for lenght
	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(nil).Select(bson.M{"track": bson.M{"$elemMatch": bson.M{"tracksrcurl": url}}}).All(&idT)
	log.Println(len(idT))
	if err != nil {
		return false

	}
	if len(idT) != 0 {
		return false
	}
	return true

}

// GetAllID we searching through the db and finds all the id and list it here
func (db *MongoDB) GetAllID() []TrackID {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []collection
	var allID []TrackID

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&all)

	if err != nil {
		return []TrackID{}
	}
	for _, data := range all {
		allID = append(allID, TrackID{data.ParagliderID.ID})
	}

	return allID
}
