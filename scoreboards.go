// pdgo Scoreboards API - unified CGO implementation

package pdgo

/*
#include <stdint.h>

// Scoreboards API
void pd_scoreboards_addScore(const char* boardID, uint32_t value, void* callback);
void pd_scoreboards_getPersonalBest(const char* boardID, void* callback);
void pd_scoreboards_freeScore(void* score);
void pd_scoreboards_getScoreboards(void* callback);
void pd_scoreboards_freeBoardsList(void* boards);
void pd_scoreboards_getScores(const char* boardID, void* callback);
void pd_scoreboards_freeScoresList(void* scores);
*/
import "C"
import "unsafe"

// Scoreboards provides access to online scoreboards
type Scoreboards struct{}

func newScoreboards() *Scoreboards {
	return &Scoreboards{}
}

// PDScore represents a score entry
type PDScore struct {
	Rank   uint32
	Value  uint32
	Player string
}

// PDScoresList represents a list of scores
type PDScoresList struct {
	BoardID     string
	LastUpdated uint32
	Scores      []PDScore
}

// PDBoard represents a scoreboard
type PDBoard struct {
	BoardID string
	Name    string
}

// PDPersonalBest represents a personal best score
type PDPersonalBest struct {
	Rank  uint32
	Value uint32
}

// ScoreboardsResult callback type
type ScoreboardsResult func(result interface{}, errorMsg string)

// AddScore adds a score to a board
func (s *Scoreboards) AddScore(boardID string, value uint32, callback ScoreboardsResult) {
	cstr := make([]byte, len(boardID)+1)
	copy(cstr, boardID)
	// Note: callback handling requires C function pointer bridging
	C.pd_scoreboards_addScore((*C.char)(unsafe.Pointer(&cstr[0])), C.uint32_t(value), nil)
}

// GetPersonalBest gets personal best score
func (s *Scoreboards) GetPersonalBest(boardID string, callback ScoreboardsResult) {
	cstr := make([]byte, len(boardID)+1)
	copy(cstr, boardID)
	C.pd_scoreboards_getPersonalBest((*C.char)(unsafe.Pointer(&cstr[0])), nil)
}

// FreeScore frees a score result
func (s *Scoreboards) FreeScore(score *PDPersonalBest) {
	// Go-allocated struct, GC handles it
}

// GetScoreboards gets list of scoreboards
func (s *Scoreboards) GetScoreboards(callback ScoreboardsResult) {
	C.pd_scoreboards_getScoreboards(nil)
}

// FreeBoardsList frees boards list
func (s *Scoreboards) FreeBoardsList(boards []PDBoard) {
	// Go-allocated, GC handles it
}

// GetScores gets scores from a board
func (s *Scoreboards) GetScores(boardID string, callback ScoreboardsResult) {
	cstr := make([]byte, len(boardID)+1)
	copy(cstr, boardID)
	C.pd_scoreboards_getScores((*C.char)(unsafe.Pointer(&cstr[0])), nil)
}

// FreeScoresList frees scores list
func (s *Scoreboards) FreeScoresList(scores *PDScoresList) {
	// Go-allocated, GC handles it
}
