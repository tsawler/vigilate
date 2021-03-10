package dbrepo

import (
	"github.com/tsawler/vigilate/internal/models"
)

// AllUsers returns all users
func (m *testDBRepo) AllUsers() ([]*models.User, error) {
	var users []*models.User

	return users, nil
}

// GetUserById returns a user by id
func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User

	return u, nil
}

// Authenticate authenticates
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// InsertRememberMeToken inserts a remember me token into remember_tokens for a user
func (m *testDBRepo) InsertRememberMeToken(id int, token string) error {
	return nil
}

// DeleteToken deletes a remember me token
func (m *testDBRepo) DeleteToken(token string) error {
	return nil
}

// CheckForToken checks for a valid remember me token
func (m *testDBRepo) CheckForToken(id int, token string) bool {
	return true
}

// Insert method to add a new record to the users table.
func (m *testDBRepo) InsertUser(u models.User) (int, error) {
	return 2, nil
}

// UpdateUser updates a user by id
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// DeleteUser sets a user to deleted by populating deleted_at value
func (m *testDBRepo) DeleteUser(id int) error {
	return nil
}

// UpdatePassword resets a password
func (m *testDBRepo) UpdatePassword(id int, newPassword string) error {
	return nil
}

// AllPreferences returns a slice of preferences
func (m *testDBRepo) AllPreferences() ([]models.Preference, error) {
	var preferences []models.Preference

	return preferences, nil
}

// SetSystemPref updates a system preference setting
func (m *testDBRepo) SetSystemPref(name, value string) error {
	return nil
}

// UpdateSystemPref updates a system preference setting
func (m *testDBRepo) UpdateSystemPref(name, value string) error {
	return nil
}

// InsertOrUpdateSitePreferences inserts or updates all site prefs from map
func (m *testDBRepo) InsertOrUpdateSitePreferences(pm map[string]string) error {
	return nil
}

// InsertHost inserts a host into the database
func (m *testDBRepo) InsertHost(h models.Host) (int, error) {
	return 2, nil
}

// GetHostByID gets a host by id and returns models.Host
func (m *testDBRepo) GetHostByID(id int) (models.Host, error) {
	var h models.Host
	return h, nil
}

// UpdateHost updates a host in the database
func (m *testDBRepo) UpdateHost(h models.Host) error {
	return nil
}

func (m *testDBRepo) GetAllServiceStatusCounts() (int, int, int, int, error) {
	return 0, 0, 0, 0, nil
}

// AllHosts returns a slice of hosts
func (m *testDBRepo) AllHosts() ([]models.Host, error) {
	var hosts []models.Host
	return hosts, nil
}

// UpdateHostServiceStatus updates the active status of a host service
func (m *testDBRepo) UpdateHostServiceStatus(hostID, serviceID, active int) error {
	return nil
}

// UpdateHostService updates a host service in the database
func (m *testDBRepo) UpdateHostService(hs models.HostService) error {
	return nil
}

// GetServicesByStatus returns all active services with a given status
func (m *testDBRepo) GetServicesByStatus(status string) ([]models.HostService, error) {
	var services []models.HostService
	return services, nil
}

// GetHostServiceByID gets a host service by id
func (m *testDBRepo) GetHostServiceByID(id int) (models.HostService, error) {
	var hs models.HostService
	return hs, nil
}

// GetServicesToMonitor gets all host services we want to monitor
func (m *testDBRepo) GetServicesToMonitor() ([]models.HostService, error) {
	var services []models.HostService
	return services, nil
}

// GetHostServiceByHostIDServiceID gets a host service by host id and service id
func (m *testDBRepo) GetHostServiceByHostIDServiceID(hostID, serviceID int) (models.HostService, error) {
	var hs models.HostService
	return hs, nil
}

// InsertEvent inserts an event into the database
func (m *testDBRepo) InsertEvent(e models.Event) error {
	return nil
}

// GetAllEvents gets all events
func (m *testDBRepo) GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	return events, nil
}
