package models

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gitlab.com/voip-services/go-kamailio-api/internal/jsonrpc"
	"gitlab.com/voip-services/go-kamailio-api/internal/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/romana/rlog"
)

const (
	subscribersTable = "subscriber"
	illegal          = "[~|{}\\[\\]<>#^’&@`]" // not allowed chars in username@domain
)

// Subscribers model
type Subscribers struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Password string `json:"password"`
	ha1      string //`json:"ha1,omitempty"`
	ha1b     string //`json:"ha1b,omitempty"`
}

// SubscribersContact is exported
type SubscribersContact struct {
	Address       string `json:"Address"`
	Expires       string `json:"Expires"`
	Q             string `json:"Q"`
	CallID        string `json:"Call-ID"`
	Cseq          string `json:"CSeq"`
	UserAgent     string `json:"User-Agent"`
	Received      string `json:"Received"`
	Path          string `json:"Path"`
	State         string `json:"State"`
	Flags         int    `json:"Falgs"`
	CFlags        int    `json:"CFlags"`
	Socket        string `json:"Socket"`
	Methods       string `json:"Methods"`
	Ruid          string `json:"Ruid"`
	Instance      string `json:"Instance"`
	RegID         int    `json:"Reg-Id"`
	ServerID      string `json:"Server-Id"`
	TcpconnID     string `json:"Tcpconn-Id"`
	Keepalive     int    `json:"Keepalive"`
	LastKeepalive int    `json:"Last-Keepalive"`
	KaRoundtrip   int    `json:"KA-Roundtrip"`
	LastModified  int    `json:"Last-Modified"`
}

// SubscribersContacts is exported
type SubscribersContacts struct {
	Contacts SubscribersContact `json:"Contact"`
}

// SubscribersInfo is exported
type SubscribersInfo struct {
	Aor      int                   `json:"AoR"`
	HashID   uint64                `json:"HashID"`
	Contacts []SubscribersContacts `json:"Contacts"`
}

// SubscribersAors is exported
type SubscribersAors struct {
	Info SubscribersInfo `json:"Info"`
}

// SubscribersStats is exported
type SubscribersStats struct {
	Records  int `json:"Records"`
	MaxSlots int `json:"Max-Slots"`
}

// SubscribersDomain is exported
type SubscribersDomain struct {
	Domain string            `json:"Domain"`
	Size   int               `json:"Size"`
	Aors   []SubscribersAors `json:"AoRs"`
	Stats  SubscribersStats  `json:"Stats"`
}

// SubscribersDomains is exported
type SubscribersDomains struct {
	Domain SubscribersDomain `json:"Domain"`
}

// SubscribersOnline is exported
type SubscribersOnline struct {
	Domains []SubscribersDomains `json:"Domains"`
}

// Prepare is exportred, do something here
func (s *Subscribers) Prepare() {
	s.Username = strings.TrimSpace(s.Username)
	s.Password = strings.TrimSpace(s.Password)
	s.Domain = strings.TrimSpace(s.Domain)
}

// Validate new subscriber. validation, ha1, ha1b calc.
func (s *Subscribers) Validate() error {

	if s.Password == "" {
		return errors.New("Password is required")
	}
	if s.Domain == "" {
		return errors.New("Domain is required")
	}
	if s.Username == "" {
		return errors.New("Username is required")
	}

	parts := strings.Split(s.Username, "@")
	re := regexp.MustCompile(illegal)

	if len(parts) > 2 { // two or more '@'
		goto regex
	} else if len(parts) == 2 { // username@domain
		if re.Match([]byte(parts[0])) || re.Match([]byte(parts[1])) {
			goto regex
		} else {
			// calculate ha1b
			var ha1b [16]byte
			ha1b = md5.Sum([]byte(parts[0] + "@" + parts[1] + ":" + s.Domain + ":" + s.Password))
			s.ha1b = hex.EncodeToString(ha1b[:])
			log.Debug("ha1b is ", s.ha1b)
			return nil
		}
	} else { // username
		if !re.Match([]byte(s.Username)) {
			// calculate ha1
			var ha1 [16]byte
			ha1 = md5.Sum([]byte(s.Username + ":" + s.Domain + ":" + s.Password))
			s.ha1 = hex.EncodeToString(ha1[:])
			log.Debug("ha1 is ", s.ha1)
			return nil
		}
		goto regex
	}

regex:
	return errors.New("Username contains illegal chars")
}

// GetSubscribersOnline is exported // TODO: mock test
func GetSubscribersOnline(httpAddr string, httpClient *http.Client) (*interface{}, error) {

	id := utils.GenerateUUID().String()
	rec := jsonrpc.NewRequest(id, "ul.dump")

	buf, err := rec.Buffer()
	if err != nil {
		log.Errorf("Failed to marshall json [%v]", err)
		return nil, err
	}
	//log.Debugf("Send: [%v]", buf)

	res, err := httpClient.Post(httpAddr, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resp := jsonrpc.Response{}
	resp.Parse(x)

	if resp.IsError() {
		return nil, errors.New(resp.Error.Error())
	}

	return &resp.Result, nil
}

// GetSubscribers func return all subscribers.
func GetSubscribers(ctx context.Context, pool *pgxpool.Pool, offset int, limit int) (*[]Subscribers, error) {

	var subs = []Subscribers{}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := `SELECT
				id,
				username,
				domain,
				password
				FROM ` + subscribersTable +
		" LIMIT $1 OFFSET $2"

	rows, _ := conn.Query(ctx, query, limit, offset)

	defer rows.Close()

	for rows.Next() {
		sub := Subscribers{}
		err := rows.Scan(
			&sub.ID,
			&sub.Username,
			&sub.Domain,
			&sub.Password)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &subs, nil
}

// GetSubscriberByID is
func GetSubscriberByID(ctx context.Context, pool *pgxpool.Pool, subID int) (*Subscribers, error) {

	var sub = Subscribers{}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	row := conn.QueryRow(ctx,
		`SELECT
			id,
			username,
			domain,
			password
			FROM `+subscribersTable+`
			WHERE id=$1`,
		subID)

	err = row.Scan(
		&sub.ID,
		&sub.Username,
		&sub.Domain,
		&sub.Password)

	// pgx.ErrNoRows
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

// GetSubscriberByUserName returns nil if subscriber not found.
func GetSubscriberByUserName(ctx context.Context, pool *pgxpool.Pool, subUserName string) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return errors.New("can not connect to the database")
	}
	defer conn.Release()

	var id Subscribers
	err = conn.QueryRow(ctx, "SELECT username FROM "+subscribersTable+" WHERE username=$1", subUserName).Scan(&id.Username)

	if err != nil {
		if err == pgx.ErrNoRows { // no username found, we are good to go
			return nil
		}
		log.Errorf("GetSubscriberByUserName error: [%v]", err)
		return err
	}

	// user already exists
	return errors.New("sip device " + id.Username + " already exists, please choose another name")

}

// Save saves subs into db
func (s *Subscribers) Save(ctx context.Context, pool *pgxpool.Pool) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	ctag, err := conn.Exec(ctx,
		"INSERT INTO "+subscribersTable+`
		(username,
		domain,
		password,
		ha1,
		ha1b)
		VALUES
		($1, $2, $3, $4, $5)`,
		s.Username,
		s.Domain,
		s.Password,
		s.ha1,
		s.ha1b)

	if err != nil {
		return err
	}
	if ctag.Insert() != true {
		log.Criticalf("Unable to INSERT: %v\n", err)
		return err
	}

	return err
}

// DeleteSubscriber removes sub from db
func DeleteSubscriber(ctx context.Context, pool *pgxpool.Pool, subid int) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	ctag, err := conn.Exec(ctx, "DELETE FROM "+subscribersTable+" WHERE id=$1", subid)
	if ctag.Delete() != true {
		log.Criticalf("Unable to DELETE: %v\n", err)
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

// UpdateSubscriber update sub
func (s *Subscribers) UpdateSubscriber(ctx context.Context, pool *pgxpool.Pool, subID int) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	rows, err := conn.Exec(ctx, `UPDATE `+subscribersTable+` SET
		username = $2,
		domain = $3,
		password = $4,
		ha1 = $5,
		ha1b = $6
		WHERE id=$1`, subID,
		s.Username,
		s.Domain,
		s.Password,
		s.ha1,
		s.ha1b)

	if rows.Update() != true {
		log.Criticalf("Unable to UPDATE: %v\n", err)
		return err
	}

	if rows.RowsAffected() > 1 {
		log.Debugf("Updated more than 1 rows [%i]", rows.RowsAffected())
	}

	if err != nil {
		return err
	}

	return nil
}
