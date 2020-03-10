package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
        "regexp"
	sys "syscall"
	"github.com/phil-fly/log"
)

var counter = 0

// Digest is a struct used for digest authentication.
// The "realm", and "nonce" fields are supplied by the server
// (in a "401 Unauthorized" response).
// The "username" and "password" fields are supplied by the client.
type Digest struct {
	Realm    string
	Nonce    string
	Username string
	Password string
}

// NewDigest returns a pointer to a new instance of authorization digest
func NewDigest() *Digest {
	return &Digest{}
}

// RandomNonce returns a random nonce
func (d *Digest) RandomNonce() {
	var timeNow sys.Timeval
	sys.Gettimeofday(&timeNow)

	counter++
	seedData := fmt.Sprintf("%d.%06d%d", timeNow.Sec, timeNow.Usec, counter)

	// Use MD5 to compute a 'random' nonce from this seed data:
	h := md5.New()
	io.WriteString(h, seedData)
	d.Nonce = hex.EncodeToString(h.Sum(nil))
}

// ComputeResponse represents generating the response using cmd and url value
func (d *Digest) ComputeResponse(cmd, url string) string {
	ha1Data := fmt.Sprintf("%s:%s:%s", d.Username, d.Realm, d.Password)
	ha2Data := fmt.Sprintf("%s:%s", cmd, url)

	h1 := md5.New()
	h2 := md5.New()
	io.WriteString(h1, ha1Data)
	io.WriteString(h2, ha2Data)

	digestData := fmt.Sprintf("%s:%s:%s", hex.EncodeToString(h1.Sum(nil)), d.Nonce, hex.EncodeToString(h2.Sum(nil)))

	h3 := md5.New()
	io.WriteString(h3, digestData)

	return hex.EncodeToString(h3.Sum(nil))
}

// AuthorizationHeader is a struct stored the infomation of parsing "Authorization:" line
type AuthorizationHeader struct {
	URI      string
	Realm    string
	Nonce    string
	Username string
	Response string
}

// ParseAuthorizationHeader represents the parsing of "Authorization:" line,
// Authorization Header contains uri, realm, nonce, Username, response fields
func ParseAuthorizationHeader(buf string) *AuthorizationHeader {
	if buf == "" {
		log.Warn("ParseAuthorizationHeader input buf is nil")
		return nil
	}

	// First, find "Authorization:"
	index := strings.Index(buf, "Authorization: Digest")
	if -1 == index {
		log.Warn("ParseAuthorizationHeader not found Authorization: Digest")
		return nil
	}

	var username, realm, nonce, uri, response string
	r := regexp.MustCompile(`username="(?s:(.*?))"`)
	text :=r.FindAllStringSubmatch(buf, -1)
	if text != nil {
		username = text[0][1]
	}

	r = regexp.MustCompile(`realm="(?s:(.*?))"`)
	text =r.FindAllStringSubmatch(buf, -1)
	if text != nil {
		realm = text[0][1]
	}

	r = regexp.MustCompile(`nonce="(?s:(.*?))"`)
	text =r.FindAllStringSubmatch(buf, -1)
	if text != nil {
		nonce = text[0][1]
	}
	r = regexp.MustCompile(`uri="(?s:(.*?))"`)
	text =r.FindAllStringSubmatch(buf, -1)
	if text != nil {
		uri = text[0][1]
	}
	r = regexp.MustCompile(`response="(?s:(.*?))"`)
	text =r.FindAllStringSubmatch(buf, -1)
	if text != nil {
		response = text[0][1]
	}
	
	log.Info("username=[%s],URI=[%s],Realm=[%s],Nonce=[%s],Response=[%s]",username,uri,realm,nonce,response)
	return &AuthorizationHeader{
		URI:      uri,
		Realm:    realm,
		Nonce:    nonce,
		Username: username,
		Response: response,
	}
}
