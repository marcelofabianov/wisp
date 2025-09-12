package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type PreferencesSuite struct {
	suite.Suite
}

func TestPreferencesSuite(t *testing.T) {
	suite.Run(t, new(PreferencesSuite))
}

func (s *PreferencesSuite) TestNewPreferences() {
	s.Run("should create preferences from a map", func() {
		data := map[string]any{"theme": "dark", "notifications": true}
		prefs, err := wisp.NewPreferences(data)
		s.Require().NoError(err)
		s.False(prefs.IsZero())

		theme, ok := prefs.Get("theme")
		s.True(ok)
		s.Equal("dark", theme)
	})

	s.Run("should create empty preferences from a nil map", func() {
		prefs, err := wisp.NewPreferences(nil)
		s.Require().NoError(err)
		s.True(prefs.IsZero())
	})

	s.Run("should parse preferences from valid JSON", func() {
		jsonData := []byte(`{"theme": "dark", "notifications": true}`)
		prefs, err := wisp.ParsePreferences(jsonData)
		s.Require().NoError(err)

		theme, ok := prefs.Get("theme")
		s.True(ok)
		s.Equal("dark", theme)
	})

	s.Run("should parse empty preferences from an empty byte slice", func() {
		prefs, err := wisp.ParsePreferences([]byte{})
		s.Require().NoError(err)
		s.True(prefs.IsZero())
	})

	s.Run("should parse empty preferences from a null JSON", func() {
		prefs, err := wisp.ParsePreferences([]byte("null"))
		s.Require().NoError(err)
		s.True(prefs.IsZero())
	})

	s.Run("should fail to parse invalid JSON", func() {
		jsonData := []byte(`{"theme": "dark"`) // JSON inválido
		_, err := wisp.ParsePreferences(jsonData)
		s.Require().Error(err)
	})
}

func (s *PreferencesSuite) TestPreferences_Immutability() {
	data := map[string]any{"lang": "pt-br"}
	prefs1, _ := wisp.NewPreferences(data)

	prefs2 := prefs1.Set("theme", "dark")

	lang, ok := prefs1.Get("lang")
	s.True(ok)
	s.Equal("pt-br", lang)
	_, ok = prefs1.Get("theme")
	s.False(ok, "prefs1 should not have the new key")

	lang, ok = prefs2.Get("lang")
	s.True(ok)
	theme, ok := prefs2.Get("theme")
	s.True(ok)
	s.Equal("dark", theme)
	s.Equal("pt-br", lang)
}

func (s *PreferencesSuite) TestPreferences_JSON_SQL() {
	prefs, _ := wisp.NewPreferences(map[string]any{"show_tutorials": false})

	s.Run("JSON Marshaling", func() {
		data, err := json.Marshal(prefs)
		s.Require().NoError(err)
		s.JSONEq(`{"show_tutorials": false}`, string(data))

		var unmarshaledPrefs wisp.Preferences
		err = json.Unmarshal(data, &unmarshaledPrefs)
		s.Require().NoError(err)
		s.Equal(prefs.Data(), unmarshaledPrefs.Data())
	})

	s.Run("SQL Interface", func() {
		val, err := prefs.Value()
		s.Require().NoError(err)
		s.IsType([]byte{}, val)

		var scannedPrefs wisp.Preferences
		err = scannedPrefs.Scan(val)
		s.Require().NoError(err)
		s.Equal(prefs.Data(), scannedPrefs.Data())
	})
}
