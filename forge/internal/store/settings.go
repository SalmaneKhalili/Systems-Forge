package store

// Settings is a tiny KV table for GUI preferences (theme, pomodoro lengths…).
func (s *Store) GetSetting(key string) (string, error) {
	var v string
	err := s.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&v)
	if err != nil {
		if isNoRows(err) {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

// SetSetting upserts a setting value.
func (s *Store) SetSetting(key, value string) error {
	_, err := s.db.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	return err
}

// AllSettings returns the full settings map.
func (s *Store) AllSettings() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

// SetSettingsBulk upserts several settings at once.
func (s *Store) SetSettingsBulk(kv map[string]string) error {
	for k, v := range kv {
		if err := s.SetSetting(k, v); err != nil {
			return err
		}
	}
	return nil
}