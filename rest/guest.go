package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Guest describes a guest record from the rest api
type Guest struct {
	Address             string                  `json:"address,omitempty"`
	ADGroup             string                  `json:"ADGroup,omitempty"`
	AgentInstalled      bool                    `json:"agentInstalled"`
	AgentMetadata       *GuestAgentMetadata     `json:"agentMetadata"`
	AgentVersion        string                  `json:"agentVersion,omitempty"`
	Backup              *GuestBackup            `json:"backup,omitempty"`
	Cpus                int                     `json:"cpus,omitempty"`
	Currentmem          int                     `json:"currentmem,omitempty"`
	Disks               []GuestDisk             `json:"disks,omitempty"`
	DuplicateMitigation bool                    `json:"duplicateMitigation"`
	Error               *GuestError             `json:"error,omitempty"`
	External            bool                    `json:"external"`
	Gateway             string                  `json:"gateway,omitempty"`
	GPU                 bool                    `json:"gpu"`
	GuestState          string                  `json:"guestState,omitempty"`
	HasBeenReady        bool                    `json:"hasBeenReady,omitempty"`
	Hostid              string                  `json:"hostid,omitempty"`
	HostDevices         []HostDevice            `json:"hostDevices,omitempty"`
	Hostname            string                  `json:"hostname,omitempty"`
	Interfaces          []GuestNetwork          `json:"interfaces,omitempty"`
	Memory              int                     `json:"memory,omitempty"`
	MigrationMetadata   *GuestMigrationMetadata `json:"migrationMetadata,omitempty"`
	MigrationProcessing bool                    `json:"migrationProcessing"`
	Name                string                  `json:"name,omitempty"`
	Os                  string                  `json:"os,omitempty"`
	Persistent          bool                    `json:"persistent,omitempty"`
	PoolID              string                  `json:"poolId,omitempty"`
	PreviousGuestState  string                  `json:"previousGuestState,omitempty"`
	ProfileID           string                  `json:"profileId,omitempty"`
	PublishedIP         string                  `json:"publishedIp,omitempty"`
	RdpUserInjected     bool                    `json:"rdpUserInjected,omitempty"`
	Realm               string                  `json:"realm,omitempty"`
	SessionInfo         *GuestSessionInfo       `json:"sessionInfo,omitempty"`
	Stamp               interface{}             `json:"stamp,omitempty"`
	Standalone          bool                    `json:"standalone"`
	StateChronology     *StateChronology        `json:"stateChronology,omitempty"`
	Tags                []string                `json:"tags,omitempty"`
	TargetState         []string                `json:"targetState,omitempty"`
	TemplateName        string                  `json:"templateName,omitempty"`
	UserVolume          *UserVolume             `json:"userVolume,omitempty"`
	Username            string                  `json:"username"`
	UserSessionState    string                  `json:"userSessionState,omitempty"`
	UserSession         *UserSession            `json:"userSession,omitempty"`
	UUID                string                  `json:"uuid,omitempty"`
	BrokerOptions       *GuestBrokerOptions     `json:"brokerOptions,omitempty"`
}

// GuestBrokerOptions allows configuring broker connection settings for a guest
type GuestBrokerOptions struct {
	Enabled           bool                    `json:"enabled"`
	DefaultConnection string                  `json:"defaultConnection,omitempty"`
	Connections       []GuestBrokerConnection `json:"connections,omitempty"`
}

// GuestBrokerConnection contains settings for a broker connection to a guest
type GuestBrokerConnection struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Port        uint   `json:"port,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Gateway     struct {
		Disabled   bool   `json:"enabled"`
		Protocols  string `json:"prococols"`
		Persistent bool   `json:"persistent"`
	} `json:"gateway"`
}

// GuestDisk is the structure for GuestDisk object in db
type GuestDisk struct {
	Type         string `json:"type"`
	Path         string `json:"path"`
	DiskDriver   string `json:"diskDriver"`
	Format       string `json:"format"`
	Filename     string `json:"filename"`
	StorageID    string `json:"storageId"`
	Size         int    `json:"size"`
	Device       string `json:"dev"`
	Backing      string `json:"backing"`
	SerialNumber string `json:"serial"`
	OSVolume     int    `json:"osvolume"`
	BootOrder    int    `json:"bootOrder"`
	Cache        string `json:"cache"`
}

// GuestNetwork is the structure for guest network interfaces
type GuestNetwork struct {
	Emulation   string `json:"emulation"`
	MacAddress  string `json:"macAddress"`
	NetworkType string `json:"network"`
	Vlan        int    `json:"vlan"`
	Bus         string `json:"bus"`
	IPAddress   string `json:"ipAddress"`
}

// GuestError is a struct for errors in the guest record
type GuestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GuestAgentMetadata contains information fabout the hive agent on the guest
type GuestAgentMetadata struct {
	ActualVersion   string `json:"actualVersion"`
	ExpectedVersion string `json:"expectedVersion"`
	Installed       bool   `json:"installed"`
	State           string `json:"state"`
	UpdateStatus    string `json:"updateStatus"`
}

// GuestSessionInfo is struct for sessioninfo of user on VM
type GuestSessionInfo struct {
	SessionID              int    `json:"sessionId"`
	SessionState           string `json:"sessionState"`
	SourceIP               string `json:"sourceIP"`
	SourceName             string `json:"sourceName"`
	SessionWorkstationName string `json:"sessionWorkstationName"`
	SessionResolution      string `json:"sessionResolution"`
}

// GuestSnapshot is contains snapshot information about the guest
type GuestSnapshot struct {
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	Current  string `json:"current"`
	State    string `json:"state"`
	Location string `json:"location"`
	Metadata string `json:"metadata"`
}

// GuestBackup is struct for backup data of VM
type GuestBackup struct {
	Enabled         bool        `json:"enabled"`
	State           string      `json:"state"`
	Frequency       string      `json:"frequency"`
	TargetStorageID string      `json:"targetStorageId"`
	LastBackup      interface{} `json:"lastBackup"`
	StateMessage    string      `json:"stateMessage"`
	DiskFrozen      bool        `json:"diskFrozen"`
}

type GuestMigrationMetadata struct {
	SourceHostId        string `json:"sourceHostId"`
	SourceHostname      string `json:"sourceHostname"`
	DestinationHostId   string `json:"destinationHostId"`
	DestinationHostname string `json:"destinationHostname"`
	Progress            int    `json:"progress"`
	MigratableXml       string `json:"migratableXml"`
}

// HostDevice is information about a device forwarded to the guest from the host
type HostDevice struct {
	Type    string `json:"type"`
	Model   string `json:"model"`
	Managed bool   `json:"managed"`
	Domain  int    `json:"domain"`
	Bus     int    `json:"bus"`
	Slot    int    `json:"slot"`
	Func    int    `json:"func"`
	UUID    string `json:"uuid"`
}

// StateChronology state tracking for the guest
type StateChronology struct {
	Current string `json:"current"`
	Next    string `json:"next"`
	Target  string `json:"target"`
}

// UserSession contains user session information for a guest
type UserSession struct {
	UserSessionState  string      `json:"userSessionState"`
	LastUserLoginTime interface{} `json:"lastUserLoginTime"`
	LastLoginDuration int         `json:"lastLoginDuration"`
	DisconnectTime    interface{} `json:"disconnectTime"`
}

// UserVolume is the struct for uservolume on guest record
type UserVolume struct {
	State                 string      `json:"state,omitempty"`
	LastReplication       interface{} `json:"lastReplication,omitempty"`
	RepliacationRequested interface{} `json:"replicationRequested,omitempty"`
	Source                string      `json:"source,omitempty"`
	Target                string      `json:"target,omitempty"`
	RunningBackup         bool        `json:"runningBackup,omitempty"`
	StateMessage          interface{} `json:"stateMessage,omitempty"`
	DiskFrozen            bool        `json:"diskFrozen,omitempty"`
	BackupSchedule        int         `json:"backupSchedule,omitempty"`
}

func (guest Guest) String() string {
	json, _ := json.MarshalIndent(guest, "", "  ")
	return string(json)
}

// ListGuests returns an array of all guests with an optional filter string
func (client *Client) ListGuests(query string) ([]Guest, error) {
	var guests []Guest
	path := "guests"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return guests, err
	}
	err = json.Unmarshal(body, &guests)
	return guests, err
}

// GetGuest requests a single guest by name
func (client *Client) GetGuest(name string) (*Guest, error) {
	var guest Guest
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	body, err := client.request("GET", "guest/"+url.PathEscape(name), nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &guest)
	return &guest, err
}

// Shutdown asks the guest operation system to shutdown
func (guest *Guest) Shutdown(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/shutdown", nil)
	return err
}

// Reboot asks the guest operation system to reboot
func (guest *Guest) Reboot(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/reboot", nil)
	return err
}

// Refresh recreates the guest with the latest pool configuration
func (guest *Guest) Refresh(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	if guest.Standalone {
		return guest.Delete(client)
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/refresh", nil)
	return err
}

// Poweron starts a powered off guest
func (guest *Guest) Poweron(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/poweron", nil)
	return err
}

// Poweroff forces a guest to powers off
func (guest *Guest) Poweroff(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/poweroff", nil)
	return err
}

// Reset forces a guest to hard reset
func (guest *Guest) Reset(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/reset", nil)
	return err
}

// Update a guest record
func (guest *Guest) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(guest)
	body, err := client.request("PUT", "guest/"+url.PathEscape(guest.Name), jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Delete deletes a guest
func (guest *Guest) Delete(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	if guest.External {
		_, err := client.request("DELETE", "guest/"+url.PathEscape(guest.Name), nil)
		return err
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/delete", nil)
	return err
}

// StartBackup requests starting a backup immediately
func (guest *Guest) StartBackup(client *Client, storageId string) (*Task, error) {
	if guest.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	data := map[string]interface{}{}
	if storageId != "" {
		data["storageId"] = storageId
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/backup", jsonValue))
}

// StartBackup requests starting a backup immediately
func (guest *Guest) ListBackups(client *Client, storageId string) ([]string, error) {
	if guest.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	data := map[string]interface{}{}
	if storageId != "" {
		data["storageId"] = storageId
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	backups := []string{}
	body, err := client.request("GET", "guest/"+url.PathEscape(guest.Name)+"/backups", jsonValue)
	if err != nil {
		return backups, err
	}
	err = json.Unmarshal(body, &backups)
	return backups, err
}

// Restore restores a guest from a backup
func (guest *Guest) Restore(client *Client, storageId, backup string) (*Task, error) {
	if guest.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	data := map[string]interface{}{}
	if storageId != "" {
		data["storageId"] = storageId
	}
	if backup != "" {
		data["backup"] = backup
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/restore", jsonValue))
}

// Migrate migrate a guest to a different host
func (guest *Guest) Migrate(client *Client, destinationHostid string) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	jsonData := map[string]string{"destinationId": destinationHostid}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/migrate", jsonValue)
	return err
}

func checkGuestState(guest Guest) bool {
	for _, v := range guest.TargetState {
		if v == guest.GuestState {
			return true
		}
	}
	return false
}

// WaitForGuest waits for a guest state to match the targetState
func (guest Guest) WaitForGuest(client *Client, timeout time.Duration) error {
	return guest.WaitForGuestWithContext(context.Background(), client, timeout)
}

// WaitForGuestWithContext waits for a guest state to match the targetState
func (guest Guest) WaitForGuestWithContext(ctx context.Context, client *Client, timeout time.Duration) error {
	if checkGuestState(guest) {
		return nil
	}
	newVal := Guest{}
	feed, err := client.GetChangeFeedWithContext(ctx, "guest", map[string]string{"name": guest.Name}, false)
	if err != nil {
		return err
	}
	timer := time.NewTimer(timeout)
	if timeout <= 0 && !timer.Stop() {
		<-timer.C
	}
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timed out")
		case msg := <-feed.Data:
			if msg.Error != nil {
				feed.Close()
				return msg.Error
			}
			err = json.Unmarshal(msg.NewValue, &newVal)
			if err != nil {
				err = fmt.Errorf("error with json unmarshal: %v", err)
				feed.Close()
				return err
			}
			if checkGuestState(newVal) {
				feed.Close()
				return nil
			}
		}
	}
}

// ResetRecord forces a guest record to be reset so it can be rebuilt
func (guest *Guest) ResetRecord(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+url.PathEscape(guest.Name)+"/resetRecord", nil)
	return err
}

// ExternalGuest is used to add external guest records to the system
type ExternalGuest struct {
	GuestName     string             `json:"guestName,omitempty"`
	Address       string             `json:"address,omitempty"`
	Username      string             `json:"username,omitempty"`
	ADGroup       string             `json:"ADGroup,omitempty"`
	Realm         string             `json:"realm,omitempty"`
	OS            string             `json:"os,omitempty"`
	BrokerOptions GuestBrokerOptions `json:"brokerOptions,omitempty"`
}

// Create creates a new pool
func (guest *ExternalGuest) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(guest)
	body, err := client.request("POST", "guest/external", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}
