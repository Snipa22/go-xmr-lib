package hashValidation

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RXVerifier struct {
	httpSession *http.Client
	uri         string
	currentSeed []byte
}

// NewRXVerifier returns a RXVerifier instance, which is used to wrap a connection to the RandomX verificaton daemon.
// It expects a server argument, which should be the HTTP interface including the port, but this is optional if
// the daemon is running on the default port/localhost bindings
func NewRXVerifier(server string) *RXVerifier {
	if len(server) == 0 {
		server = "http://127.0.0.1:39093"
	}
	return &RXVerifier{
		httpSession: &http.Client{},
		uri:         server,
	}
}

// Hash performs an actual round of hashing against a given input, for RandomX, this sends it off to the verifier
func (s *RXVerifier) Hash(input []byte, seed []byte) ([]byte, error) {
	if bytes.Compare(seed, s.currentSeed) != 0 {
		if err := s.NewSeed(seed); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(fmt.Sprintf("%v/seed", s.uri), "application/x.randomx+bin", bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	req.Header.Set("RandomX-Seed", hex.EncodeToString(seed))
	resp, err := s.httpSession.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		if b, err := io.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			return hex.DecodeString(string(b))
		}
	case 400:
		return nil, invalidBody
	case 403:
		return nil, notInitialized
	case 413:
		return nil, payloadSize
	case 415:
		return nil, badHeader
	case 422:
		return nil, badSeed
	default:
		return nil, nil
	}
}

// NewSeed sets the seed value into the remote daemon and updates the local seed directly.
func (s *RXVerifier) NewSeed(input []byte) error {
	if bytes.Compare(s.currentSeed, input) == 0 {
		return nil
	}
	resp, err := s.httpSession.Post(fmt.Sprintf("%v/seed", s.uri), "application/x.randomx+bin", bytes.NewReader(input))
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 204:
		s.currentSeed = input
		return nil
	case 400:
		return invalidBody
	case 413:
		return invalidBody
	case 415:
		return badHeader
	default:
		return invalidResponse
	}
}

// Info returns the information available from the /info endpoint in the daemon
func (s *RXVerifier) Info() (*ValidatorInfo, error) {
	resp, err := s.httpSession.Get(fmt.Sprintf("%v/info"))
	if err != nil {
		return nil, err
	}
	if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		retVal := &ValidatorInfo{}
		if err = json.Unmarshal(b, retVal); err != nil {
			return nil, err
		}
		return retVal, nil
	}
}

// SetCurrentSeed is used to overwrite the current seed in an RXVerifier setup, this is largely for multithreaded
// implementations where multiple clients need to synchronize their seed across instances together
func (s *RXVerifier) SetCurrentSeed(input []byte) {
	s.currentSeed = input
}

var badSeed = errors.New("seed in the hashing daemon does not match provided seed")
var notInitialized = errors.New("hashing daemon not initialized")
var invalidBody = errors.New("invalid body")
var badHeader = errors.New("bad content-type header")
var invalidResponse = errors.New("unhandled status response from verification daemon")
var payloadSize = errors.New("the POST body is too large")
