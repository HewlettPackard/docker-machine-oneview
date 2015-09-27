package icsp

import "github.com/docker/machine/log"

// ValueItem struct
type ValueItem struct {
	Scope string `json:"scope,omitempty"` // scope of value
	Value string `json:"value,omitempty"` // value of information
}

// CustomAttribute struct
type CustomAttribute struct {
	Key    string      `json:"key,omitempty"` // key for name value pairs
	Values []ValueItem `json:"values,omitempty"`
}

// GetValueItem - gets a ValueItem from custom attribute
func (s *Server) GetValueItem(key string, scope string) (int, ValueItem) {
	var v ValueItem
	i, vitems := s.GetValueItems(key)
	if i >= 0 {
		for i, vitem := range vitems {
			if vitem.Scope == scope {
				return i, vitem
			}
		}
	}
	return -1, v
}

// GetValueItems - gets a customattribute value item by key
func (s *Server) GetValueItems(key string) (int, []ValueItem) {
	var vi []ValueItem
	for i, attribute := range s.CustomAttributes {
		if attribute.Key == key {
			return i, attribute.Values
		}
	}
	return -1, vi
}

// SetValueItems object
func (s *Server) SetValueItems(key string, newv ValueItem) {
	_, oldv := s.GetValueItem(key, newv.Scope)
	log.Debugf("ValueItem => %+v", oldv)

	if i, oldv := s.GetValueItem(key, newv.Scope); i < 0 {
		// creat a new ValueItem
		log.Debugf("Adding new ValueItem => %v+", newv)
		vi, _ := s.GetValueItems(key)
		if vi < 0 {
			// a new key is needed
			s.CustomAttributes = append(s.CustomAttributes, CustomAttribute{Key: key, Values: []ValueItem{{Scope: newv.Scope, Value: newv.Value}}})
		} else {
			s.CustomAttributes[vi].Values = append(s.CustomAttributes[vi].Values, ValueItem{Scope: newv.Scope, Value: newv.Value})
		}
	} else {
		// set an existing one
		log.Debugf("Change %v+ to >>  %v+", oldv, newv)
		vi, _ := s.GetValueItems(key)
		s.CustomAttributes[vi].Values[i] = ValueItem{Scope: newv.Scope, Value: newv.Value}
	}
}

// SetCustomAttribute  set a custom attribute for server
func (s *Server) SetCustomAttribute(key string, scope string, value string) {
	s.SetValueItems(key, ValueItem{Scope: scope, Value: value})
}
