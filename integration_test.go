package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// User represents the model for the user object in the test
type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// TestIntegrationSaveAndRetrieve tests the save and retrieve endpoints of the microservice.
func TestIntegrationSaveAndRetrieve(t *testing.T) {
	server := "http://localhost:8080"
	// Step 1: Create a test user with a UUID filled with 'E' hex characters
	testUser := User{
		ID:          uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:        "Jane Doe",
		Email:       "j_d@example.com",
		DateOfBirth: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	// Step 2: Now test the retrieval using the test UUID
	retrieveURL := fmt.Sprintf("%s/%s", server, testUser.ID.String())
	req, err := http.NewRequest("GET", retrieveURL, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	// Step 3: Send the GET request to retrieve the user
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Step 4: Check the status code for the /{id} endpoint
	if http.StatusNotFound != resp.StatusCode {
		output, err := deleteTestUser(testUser)
		if err != nil {
			t.Fatalf("Failed to clean up database: %v\nOutput: %s", err, output)
		}
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}
		defer resp.Body.Close()
	}
	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected status 404 OK from /{id} endpoint")

	// Step 5: Marshal the testUser to JSON
	userJson, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("Failed to marshal test user: %v", err)
	}

	// Step 6: Create a request to the /save endpoint
	saveURL := fmt.Sprintf("%s/save", server)
	req, err = http.NewRequest("POST", saveURL, bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Step 7: Send the request and check the response
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Step 8: Check the response status code for the /save endpoint
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 200 OK from /save endpoint")

	// Step 9: Now test the retrieval using the same UUID
	retrieveURL = fmt.Sprintf("%s/%s", server, testUser.ID.String())
	req, err = http.NewRequest("GET", retrieveURL, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	// Step 10: Send the GET request to retrieve the user
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Step 11: Read the response body using io.ReadAll (instead of ioutil.ReadAll)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Step 12: Check the status code for the /{id} endpoint
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK from /{id} endpoint")

	// Step 13: Unmarshal the response body back to a User object
	var retrievedUser User
	err = json.Unmarshal(body, &retrievedUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Step 14: Assert that the retrieved user matches the saved user
	assert.Equal(t, testUser.ID, retrievedUser.ID, "Retrieved user ID should match the saved user ID")
	assert.Equal(t, testUser.Name, retrievedUser.Name, "Retrieved user name should match the saved user name")
	assert.Equal(t, testUser.Email, retrievedUser.Email, "Retrieved user email should match the saved user email")
	assert.Equal(t, testUser.DateOfBirth, retrievedUser.DateOfBirth, "Retrieved user date of birth should match the saved user date of birth")

	// Step 15: Delete test user inside the Docker container
	output, err := deleteTestUser(testUser)
	if err != nil {
		t.Fatalf("Failed to clean up database: %v\nOutput: %s", err, output)
	}
}

func deleteTestUser(testUser User) (string, error) {
	deleteCmd := fmt.Sprintf("psql -U myuser -d mydb -c \"DELETE FROM users WHERE id='%s';\"", testUser.ID.String())
	cmd := exec.Command("docker", "exec", "user_db-database", "sh", "-c", deleteCmd)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
