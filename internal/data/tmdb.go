package data

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/danmharris/random-episode/internal/episode"
)

type TMDB struct {
	client *tmdb.Client
}

func NewTMDB(bearerToken string) (*TMDB, error) {
	client, err := tmdb.InitV4(bearerToken)
	if err != nil {
		return nil, err
	}

	return &TMDB{client}, nil
}

func (tmdb *TMDB) ShowDetails(id int) (*episode.Show, error) {
	details, err := tmdb.client.GetTVDetails(id, nil)
	if err != nil {
		return nil, err
	}

	seasonCount := details.NumberOfSeasons
	episodeCounts := make([]int, seasonCount)

	for i := 0; i < seasonCount; i++ {
		episodeCounts[i] = details.Seasons[i].EpisodeCount
	}

	return &episode.Show{
		ShowID:       id,
		Title:        details.Name,
		EpisodeCount: episodeCounts,
	}, nil
}

func (tmdb *TMDB) EpisodeDetails(sh *episode.Show, pos episode.EpisodePosition) (*episode.Episode, error) {
	details, err := tmdb.client.GetTVEpisodeDetails(sh.ShowID, pos.Season, pos.Episode, nil)
	if err != nil {
		return nil, err
	}

	return &episode.Episode{
		EpisodePosition: pos,
		Title:           details.Name,
		Show:            sh,
	}, nil
}
