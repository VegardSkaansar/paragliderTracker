# paragliderTracker

This is a paraglidetracker where you can post your tracking information, and this api i will help you get the trackpoints of this tracks

heroku : https://paraglidertracker.herokuapp.com/paragliding/

# Mgo pkg for mongodb
I started downloading this version and test it before i heard about globalsign/mgo from github. Since this version of mgo already was working for me i kept using the first version.

# Clocktrigger 
At first i didnt understand how i would get my clocktrigger on openstack but after the lecture of docker compose i got an idea of how to do it. But this was after the submition, so this have not
been implemented.

#Paraglider track

## GET /paragliding/api
Returns the uptime of the application.

## POST /paragliding/api/track
Registers a track, and one track can only be added once.

## GET /paragliding/api/track
Gets a array with all the unique track ids.

## GET /paragliding/api/track/<id>
Returns metadata from a track with the specified id.

## GET /paragliding/api/track/<id>/<field>
Returns a specific field from the metadata as text.

# Ticker api

## GET /paragliding/api/ticker/
Returns the timestamp of the latest added track.

##GET /paragliding/api/ticker/<timestamp>
Returns all tracks with higher timestamp than <timestamp>.

# Webhook api

## POST /paragliding/api/webhook/new_track
Registers a webhook that will be notifyed when a new track is added.

## GET /paragliding/api/webhook/new_track/<webhook_id>
Gets details of a webhook.

## Delete /paragliding/api/webhook/new_track/<webhook_id>
Delete the webhook with this specified id.
