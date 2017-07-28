package mlib

import "testing"

func TestOps(t *testing.T) {
	mm := NewMusicManager()

	if mm == nil {
		t.Error("New Music Manager failed.")
	}

	if mm.len() != 0 {
		t.Error("New Music Manager failed , not empty.")
	}

	m0 := &MusicEntry{
		"1", "My Heart Will Go On.", "Celion Dion", Pop,
		"http://gbox.me/12341234", MPS}

	mm.Add(m0)

	if m.Len() != 1 {
		t.Error("Music Manager.Add() failed.")
	}

	m := mm.Find(m0.name)
	if m == nil {
		t.Error("Music Manager.Find() failed.")
	}

	if m.Id != m0.Id || m.Artist != m.Artist {

	}

}
