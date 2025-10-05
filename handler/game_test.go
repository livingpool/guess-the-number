package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/livingpool/middleware"
	"github.com/livingpool/mocks"
	"github.com/livingpool/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGame(t *testing.T) {
	mockTemplates := mocks.NewMockTemplatesRepository(t)
	mockPlayerPool := mocks.NewMockPlayerPoolRepository(t)
	mockTimeProvider := mocks.NewMockTimeProviderRepository(t)
	mockGameHandler := NewGameHandler(mockTemplates, mockPlayerPool, mockTimeProvider)

	testcases := []struct {
		name         string
		digit        string
		expectedCode int
		expectedMsg  string
	}{
		{"alphabet_input", "a", 422, "Input is not a digit"},
		{"non_alnum_input", "!", 422, "Input is not a digit"},
		{"lower_than_range_input", "0", 422, "Input is not in range"},
		{"higher_than_range_input", "15", 422, "Input is not in range"},
		{"valid_input", "8", 200, "8"},
	}

	mockPlayerPool.EXPECT().NewPlayer(mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(new(service.Player))
	mockPlayerPool.EXPECT().AddPlayer(mock.Anything).Return(nil)

	received := make([]service.FormData, 0, len(testcases))
	mockTemplates.EXPECT().Render(
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).RunAndReturn(func(w io.Writer, tmpl string, data any) error {
		formData, ok := data.(service.FormData)
		if !ok {
			t.Errorf("data is not service.FormData")
		}
		received = append(received, formData)
		return nil
	})

	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/home?digit="+tc.digit, nil)

			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIdKey, uuid.New().String()))

			mockGameHandler.NewGame(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedCode >= 400 {
				assert.Contains(t, received[i].Error, tc.expectedMsg)
			} else {
				assert.Equal(t, strconv.Itoa(received[i].Digit), tc.expectedMsg)
			}
		})
	}
}

func TestGenHint(t *testing.T) {
	testcases := []struct {
		guess    string
		answer   string
		expected string
	}{
		{"0", "1", "0a0b"},
		{"1", "1", "1a0b"},
		{"10", "11", "1a0b"},
		{"12", "21", "0a2b"},
		{"111", "111", "3a0b"},
		{"22300", "00322", "1a4b"},
		{"123456", "654321", "0a6b"},
	}

	for _, tc := range testcases {
		actual := genHint(tc.guess, tc.answer)
		assert.Equal(t, tc.expected, actual)
	}
}

// this seems like a bit too much mocking tbf
func TestCheckGuess(t *testing.T) {
	mockTemplates := mocks.NewMockTemplatesRepository(t)
	mockPlayerPool := mocks.NewMockPlayerPoolRepository(t)
	mockTimeProvider := mocks.NewMockTimeProviderRepository(t)
	mockGameHandler := NewGameHandler(mockTemplates, mockPlayerPool, mockTimeProvider)

	testcases := []struct {
		name           string
		guess          string
		playerId       string
		expectedResult string
		expectedCode   int
		expectedMsg    string
	}{
		{"invalid_player_id", "", "a", "", 404, "player id is not an integer"},
		{"player_id_not_found", "", "-1", "", 404, "player id not found"},
		{"empty_guess", "", "1", "", 422, ""},
		{"wrong_length_guess", "12", "0", "", 422, ""},
		{"correct_guess", "111", "111", "3a0b", 200, ""},
	}

	mockTimeProvider.EXPECT().Now(mock.Anything).Return(time.Now())

	// TODO: test template?
	mockTemplates.EXPECT().Render(
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).RunAndReturn(func(w io.Writer, tmpl string, data any) error {
		return nil
	})

	// i set playerId and answer to the same value so i can mock this easier
	// so note that the testcases should follow this setting
	mockPlayerPool.EXPECT().GetPlayer(mock.AnythingOfType("int")).RunAndReturn(
		func(id int) (*service.Player, bool) {
			if id >= 0 {
				return &service.Player{Id: id, Answer: strconv.Itoa(id)}, true
			}
			return nil, false
		},
	)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/check?guess="+tc.guess+"&id="+tc.playerId, nil)

			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIdKey, uuid.New().String()))

			mockGameHandler.CheckGuess(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedCode >= 400 {
				assert.Contains(t, w.Body.String(), tc.expectedMsg)
			}
		})
	}
}
