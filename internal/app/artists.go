package app

import "github.com/bigdaddy284/mymusic/internal/media"

// FollowedArtists returns the artists the user follows in MyTube, the
// last followed first.
func (c *Collection) FollowedArtists() []media.Artist { return c.a.library.Artists() }

// FollowsArtist reports whether the user follows the artist id in
// MyTube.
func (c *Collection) FollowsArtist(id string) bool { return c.a.library.FollowsArtist(id) }

// SetArtistFollowed follows an artist in MyTube, or stops following
// them there. Nothing of it reaches their service, signed in or not:
// SetFollowing follows there.
func (c *Collection) SetArtistFollowed(a media.Artist, on bool) error {
	if err := c.a.library.SetArtistFollowed(a, on); err != nil {
		return err
	}
	if !on {
		c.a.releases.forget(c.a, a.ID)
	}
	return nil
}
