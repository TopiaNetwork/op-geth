package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type KVDB interface {
	Get(key []byte) ([]byte, error)
	Create(key, value []byte) error
	Update(key, value []byte) error
	Delete(key []byte) error
	Close() error
}

type LevelKVDB struct {
	Url string
}

func NewLevelKVDB(url string) *LevelKVDB {
	return &LevelKVDB{Url: url}
}

func (db *LevelKVDB) Get(key []byte) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%v/get/%s", db.Url, key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get request failed with status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return []byte(result["value"]), nil
}

func (db *LevelKVDB) Create(key, value []byte) error {
	data := fmt.Sprintf(`{"key":"%s", "value":"%s"}`, string(key), string(value))
	resp, err := http.Post(fmt.Sprintf("%v/put", db.Url), "application/json", strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("put request failed with status: %d", resp.StatusCode)
	}
	return nil
}

func (db *LevelKVDB) Update(key, value []byte) error {
	return nil
}

func (db *LevelKVDB) Delete(key []byte) error {
	return nil
}

func (db *LevelKVDB) Close() error {
	return nil
}
