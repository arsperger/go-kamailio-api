# go-kamailio-api

REST API over HTTP and JSONRPC client
Web server listens on HTTP and talks JSON-RPC 2.0

The main idea is to abstract back-end VoIP related configuration complexities and provide simple endpoints for provisioning and management of voice core. This repo is simple example of CRUD operations for SIP accounts  and fetching registered devices.
The name contains kamailio because initially it was developed for Kamailio based voip platform.

Golang 1.16

## Installation

### Requirements

* SIP proxy
* Database: Postgresql
* Golang

The easiest way is to build docker image, Dockerfile is in the root folder.

Alternatively you can clone the source code from this repo download the dependencies: `go get -d ./...` and build it with `make build` command.

## Configuration

Configuration is in `config.json`
Database connection string could be passed as an environment variable `KAM_DB_URL=postgres://kamailio:kamailio@localhost/kamailio` or set in the config file.
Currently supported only `Postgresql` database via `pgx` driver v4.11.0 with configurable connection pool.

path prefix `/api/v1`

### SIP proxy

TBD

## API endpoints

currently implemented only for mod `auth_db` - subscribers table.

Request | Endpoints           | Functionality
--------|---------------------|--------------------------------
POST    | /subscribers        | Create SIP device ( see example )
GET     | /subscribers        | Get all SIP devices
GET     | /subscribers/online | Get currently registered SIP devices (from internal location struct)
GET     | /subscribers/id     | Get SIP device by ID
PUT     | /subscribers/id     | Update SIP device by ID
DELETE  | /subscribers/id     | Delete SIP deivce by ID

### Examples

*Create SIP device:*

```json
curl -X POST -d {
"username": "7777",
"domain": "example.com",
"password": "sdsd",
} http://localhost:8080/api/v1/subscribers

```

*response:*

TBD

*GET SIP device by ID:*

`GET http://localhost:8080/api/v1/subscribers/15`

*response:*

```json
{"id":15,"username":"arsen","password":"sdsd","caller_name":"test1212","caller_number":"1212","active":true,"enable_push":true,"account_id":2,"agreement_id":2,"sip_profile_id":1,"allow_local_calls":false,"incoming_pricelist_id":2}
```

*Get registered SIP devices (Online users)*

`GET http://localhost:8080/api/v1/subscribers/online`

*response:*

```json
{"Domains":[{"Domain":{"AoRs":[{"Info":{"AoR":"7777","Contacts":[{"Contact":{"Address":"sip:7777@172.16.238.20:5060;ob;alias=172.16.238.20~5060~1","CFlags":4,"CSeq":55297,"Call-ID":"JWHwQwRDg38KCD88kZesPN810eTg973l","Expires":84,"Flags":1,"Instance":"[not set]","KA-Roundtrip":0,"Keepalive":1,"Last-Keepalive":1614616555,"Last-Modified":1614616555,"Methods":8159,"Path":"[not set]","Q":-1,"Received":"[not set]","Reg-Id":0,"Ruid":"uloc-603d17af-27-1","Server-Id":0,"Socket":"udp:172.16.238.14:5060","State":"CS_NEW","Tcpconn-Id":-1,"User-Agent":"PJSUA v2.9 Linux-5.9.198.212/x86_64"}}],"HashID":836344572}}],"Domain":"kamailio_location","Size":1024,"Stats":{"Max-Slots":1,"Records":1}}}]}
```

## TODO

* add more mods (dispatcher, dialplan, etc)
* more tests, httptest for API calls
* Websocket server
* ...
