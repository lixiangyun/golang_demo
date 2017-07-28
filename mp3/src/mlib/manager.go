package mlib

import "errors"

type MusicEntry struct {
	Id     string
	Name   string
	Artist string
	Source string
	Type   string
}

type MusicManager struct {
	music []MusicEntry
}

func NewMusicManager() *MusicManager {
	return &MusicManager{make([]MusicEntry, 0)}
}

func (m *MusicManager) Len() int {
	return len(m.music)
}

func (m *MusicManager) Get(index int) (music *MusicEntry, err error) {
	if index < 0 || index >= len(m.music) {
		return nil, errors.New("Index out of range.")
	}

	return &m.music[index], nil
}

func (m *MusicManager) Find(name string) *MusicEntry {
	if len(m.music) == 0 {
		return nil
	}

	for _, v := range m.music {
		if v.Name == name {
			return &v
		}
	}

	return nil
}

func (m *MusicManager) Add(music *MusicEntry) {
	m.music = append(m.music, *music)
}

func (m *MusicManager) Remove(index int) *MusicEntry {
	if index < 0 || index >= len(m.music) {
		return nil
	}
	remusic := &m.music[index]
	m.music = append(m.music[:index], m.music[index+1:]...)
	return remusic
}
