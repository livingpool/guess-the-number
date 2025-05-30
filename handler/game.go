package handler

import (
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Livingpool/constants"
	"github.com/Livingpool/middleware"
	"github.com/Livingpool/service"
	"github.com/Livingpool/utils"
	"github.com/Livingpool/views"
)

type GameHandler struct {
	renderer     views.TemplatesInterface
	playerPool   service.PlayerPoolInterface
	timeProvider service.TimeProviderInterface
}

func NewGameHandler(r views.TemplatesInterface, p service.PlayerPoolInterface, t service.TimeProviderInterface) *GameHandler {
	return &GameHandler{
		renderer:     r,
		playerPool:   p,
		timeProvider: t,
	}
}

func (h *GameHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, "base", nil)
}

func (h *GameHandler) ReturnHome(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, "home", nil)
}

func (h *GameHandler) NewGame(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(middleware.RequestIdKey).(string)
	digit, err := strconv.Atoi(r.FormValue("digit"))

	// Invalid input error
	if err != nil {
		formData := service.FormData{
			Error: "Input is not a digit :(",
		}
		w.Header().Set("HX-Retarget", "#form")
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.renderer.Render(w, "form", formData)
		return
	}

	// Input value not in range error
	if digit < constants.DIGIT_LOWER_LIMIT || digit > constants.DIGIT_UPPER_LIMIT {
		formData := service.FormData{
			Error: "Input is not in range :(",
		}
		w.Header().Set("HX-Retarget", "#form")
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.renderer.Render(w, "form", formData)
		return
	}

	lower, upper := calcRange(digit)
	answer := genRandInt(lower, upper)
	answerStr := strconv.Itoa(answer)

	// PlayerPool full error
	newPlayer := h.playerPool.NewPlayer(answerStr)
	if err = h.playerPool.AddPlayer(newPlayer); err != nil {
		formData := service.FormData{
			Error: "Server is full. Please try again later!",
		}
		w.Header().Set("HX-Retarget", "#form")
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.renderer.Render(w, "form", formData)
		return
	}

	formData := service.FormData{
		Digit:    digit,
		Start:    strconv.Itoa(lower),
		End:      strconv.Itoa(upper),
		Error:    "",
		PlayerId: newPlayer.Id,
	}

	slog.Info("player registered", "reqId", reqId, "playerId", newPlayer.Id)

	if err := h.renderer.Render(w, "game", formData); err != nil {
		slog.Error("render game error", "err", err.Error())
	}
}

func (h *GameHandler) CheckGuess(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(middleware.RequestIdKey).(string)
	guessStr := r.URL.Query().Get("guess")
	playerId := r.URL.Query().Get("id")

	// Player id not parseable error
	id, err := strconv.Atoi(playerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("player id is not an integer"))
		return
	}

	// Player doesn't exist error
	player, exists := h.playerPool.GetPlayer(id)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("player id not found"))
		slog.Error("player id not found", "reqId", reqId, "playerId", id)
		return
	}

	// Guess/Answer length don't match, but as this is a valid user input, we still return the guess results
	// at the frontend, a Swal error popup is generated for this.
	if len(guessStr) != len(player.Answer) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.renderer.Render(w, "result", player.GuessResults)
		return
	}

	result := genHint(guessStr, player.Answer)
	timeZone := utils.GetTimeZone(utils.ReadUserIP(r))

	row := service.ResultRow{
		TimeStamp: h.timeProvider.Now(timeZone).Format(time.TimeOnly),
		Guess:     "#" + strconv.Itoa(len(player.GuessResults.Rows)+1) + ": " + guessStr,
		Result:    result,
	}

	player.GuessResults.Rows = append([]service.ResultRow{row}, player.GuessResults.Rows...)

	if err := h.renderer.Render(w, "result", player.GuessResults); err != nil {
		slog.Error("render game error", "err", err.Error())
	}
}

func genHint(guess, answer string) string {
	a, b := 0, 0
	aMap := make([]bool, len(guess)) // positions of a's
	countMap := make([]int, 10)      // occurences of chars (for calc b's)

	for i := range 10 {
		countMap[i] = strings.Count(answer, strconv.Itoa(i))
	}

	// log.Println(guessStr, player.Answer)
	for i := range len(guess) {
		c, _ := strconv.Atoi(string(guess[i]))
		if guess[i] == answer[i] {
			if countMap[c] <= 0 {
				b--
				a++
			} else {
				a++
				countMap[c]--
			}
			aMap[i] = true
		} else {
			if countMap[c] > 0 {
				b++
				countMap[c]--
			}
		}
	}

	countA := strconv.Itoa(a)
	countB := strconv.Itoa(b)
	return countA + "a" + countB + "b"
}

// returns [lower, upper] from given digit
func calcRange(digit int) (int, int) {
	if digit == 1 {
		return 0, 9
	}
	upper, _ := strconv.Atoi(strings.Repeat("9", digit))
	return int(math.Pow(10, float64(digit-1))), upper
}

// generate a pseudo random number with range [lower, upper]
func genRandInt(lower, upper int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(upper-lower) + lower
	return num
}
