package episode

import "math/rand"

type Show struct {
	ShowID       int
	Title        string
	EpisodeCount []int
}

type EpisodePosition struct {
	Season  int
	Episode int
}

type Episode struct {
	EpisodePosition
	Show  *Show
	Title string
}

type EpisodeData interface {
	ShowDetails(id int) (*Show, error)
	EpisodeDetails(show *Show, epos EpisodePosition) (*Episode, error)
}

func (s Show) RandomEpisode() EpisodePosition {
	season := rand.Intn(len(s.EpisodeCount))
	episode := rand.Intn(s.EpisodeCount[season])

	return EpisodePosition{season + 1, episode + 1}
}

func FindRandomEpisode(ed EpisodeData, showID int) (*Episode, error) {
	show, err := ed.ShowDetails(showID)
	if err != nil {
		return nil, err
	}

	pos := show.RandomEpisode()
	episode, err := ed.EpisodeDetails(show, pos)
	if err != nil {
		return nil, err
	}

	return episode, nil
}
