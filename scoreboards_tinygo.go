//go:build tinygo

// TinyGo implementation of Scoreboards API

package pdgo

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
	if bridgeScoreboardsAddScore != nil {
		buf := make([]byte, len(boardID)+1)
		copy(buf, boardID)
		// Note: callback handling requires C function pointer bridging
		bridgeScoreboardsAddScore(&buf[0], value, 0)
	}
}

// GetPersonalBest gets personal best score
func (s *Scoreboards) GetPersonalBest(boardID string, callback ScoreboardsResult) {
	if bridgeScoreboardsGetPersonalBest != nil {
		buf := make([]byte, len(boardID)+1)
		copy(buf, boardID)
		bridgeScoreboardsGetPersonalBest(&buf[0], 0)
	}
}

// FreeScore frees a score result
func (s *Scoreboards) FreeScore(score *PDPersonalBest) {
	if bridgeScoreboardsFreeScore != nil && score != nil {
		// The score struct is Go-allocated, just let GC handle it
	}
}

// GetScoreboards gets list of scoreboards
func (s *Scoreboards) GetScoreboards(callback ScoreboardsResult) {
	if bridgeScoreboardsGetScoreboards != nil {
		bridgeScoreboardsGetScoreboards(0)
	}
}

// FreeBoardsList frees boards list
func (s *Scoreboards) FreeBoardsList(boards []PDBoard) {
	// Go-allocated, GC handles it
}

// GetScores gets scores from a board
func (s *Scoreboards) GetScores(boardID string, callback ScoreboardsResult) {
	if bridgeScoreboardsGetScores != nil {
		buf := make([]byte, len(boardID)+1)
		copy(buf, boardID)
		bridgeScoreboardsGetScores(&buf[0], 0)
	}
}

// FreeScoresList frees scores list
func (s *Scoreboards) FreeScoresList(scores *PDScoresList) {
	// Go-allocated, GC handles it
}
