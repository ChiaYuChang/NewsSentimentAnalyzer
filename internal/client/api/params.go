package api

import (
	"encoding/json"
	"net/url"
	"strings"
)

type Params struct {
	url.Values
	sep string
}

func NewParams(sep string) *Params {
	return &Params{Values: url.Values{}, sep: sep}
}

func (p Params) Sep() string {
	return p.sep
}

func (p Params) Add(key Key, val string) {
	if val == "" {
		return
	}
	p.Values.Add(string(key), val)
}

// set parameter to the given value, if val is empty, nothing happens
func (p Params) Set(key Key, val string) {
	if val == "" {
		return
	}
	p.Values.Set(string(key), val)
}

func (p Params) Del(key Key) {
	p.Values.Del(string(key))
}

func (p Params) Get(key Key) string {
	return p.Values.Get(string(key))
}

func (p Params) Has(key Key) bool {
	return p.Values.Has(string(key))
}

func (p Params) MarshalJSON() ([]byte, error) {
	if p.sep == "" {
		return json.Marshal(p.Values)
	}

	vals := url.Values{}
	for k, v := range p.Values {
		vals[k] = []string{strings.Join(v, p.sep)}
	}
	return json.Marshal(vals)
}

func (p *Params) UnmarshalJSON(b []byte) error {
	vals := url.Values{}
	if err := json.Unmarshal(b, &vals); err != nil {
		return err
	}

	if p.sep == "" {
		p.Values = vals
		return nil
	}

	p.Values = url.Values{}
	for k, vcsl := range vals {
		for _, v := range strings.Split(vcsl[0], p.sep) {
			p.Values.Add(k, v)
		}
	}
	return nil
}

func (p Params) Encode() string {
	if p.sep == "" {
		return p.Values.Encode()
	}

	vals := url.Values{}
	for k, v := range p.Values {
		vals[k] = []string{strings.Join(v, p.sep)}
	}
	return vals.Encode()
}

func (p *Params) Decode(q string) error {
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}

	if p.sep == "" {
		p.Values = vals
		return nil
	}

	p.Values = url.Values{}
	for k, vcsl := range vals {
		for _, v := range strings.Split(vcsl[0], p.sep) {
			p.Values.Add(k, v)
		}
	}
	return nil
}

func (p0 *Params) Clone() (*Params, error) {
	b, err := json.Marshal(p0)
	if err != nil {
		return nil, err
	}

	p1 := NewParams(p0.sep)
	err = json.Unmarshal(b, p1)
	if err != nil {
		return nil, err
	}
	return p1, nil
}
