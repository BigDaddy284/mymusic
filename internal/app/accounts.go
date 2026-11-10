package app

import (
	"context"
	"log"
	"time"
)

// loginTimeout is how long signing in waits for the user in the browser.
const loginTimeout = 5 * time.Minute

// spotifySecret is the key under which MyTube 0.x kept a Spotify sign-in
// in the system's credential store, which MyTube now forgets
// (forgetSpotify).
const spotifySecret = "spotify-token"

// AccountsStatus is the state of the user's accounts.
type AccountsStatus struct {
	YouTube YouTubeAccount `json:"youtube"`
	Lastfm  LastfmAccount  `json:"lastfm"`
}

// Accounts signs the user in and out of their music services.
type Accounts struct{ a *App }

// profileRetry is how long a profile that could not be read waits before
// it is read again. Every window asks for the accounts' state again once a
// reading tells it, so a profile read again at once would be read again
// and again, as when a service refuses the sign-in, and the windows told
// again and again (#50).
const profileRetry = time.Minute

// profileErrors is why the profiles read in the background could not be
// read, and when they were.
type profileErrors struct {
	youtube string
	at      time.Time
}

// waits reports whether a profile that could not be read, for err, waits
// before it is read again.
func (e profileErrors) waits(err string) bool {
	return err != "" && time.Since(e.at) < profileRetry
}

// Status returns the state of the user's accounts at once, without asking
// the services: whether the user is signed in, and the profiles read
// before. The main window waits for it before it shows anything, so the
// profiles not read yet are read in the background, and AccountsChanged
// brings them to every window. A profile that could not be read comes with
// why, and is read again after profileRetry.
func (s *Accounts) Status(ctx context.Context) AccountsStatus {
	return s.a.knownAccountsStatus()
}

// knownAccountsStatus is the state of the accounts as MyTube knows it,
// without asking the services; it has the profiles it lacks read in the
// background.
func (a *App) knownAccountsStatus() AccountsStatus {
	a.mu.Lock()
	youtubeMe, errs := a.youtubeMe, a.profileErrs
	a.mu.Unlock()
	status := AccountsStatus{
		YouTube: YouTubeAccount{SignedIn: a.youtube.SignedIn(), Account: youtubeMe},
		Lastfm:  a.lastfmAccount(),
	}
	if status.YouTube.SignedIn && youtubeMe == nil {
		status.YouTube.Error = errs.youtube
		if !errs.waits(errs.youtube) {
			a.readProfiles()
		}
	}
	return status
}

// readProfiles reads the profiles of the accounts in the background, one
// reading at a time, and tells every window the accounts' state once read.
func (a *App) readProfiles() {
	if !a.readingProfiles.CompareAndSwap(false, true) {
		return
	}
	go func() {
		defer a.readingProfiles.Store(false)
		// Reading the status reads the profiles, and keeps them; it says
		// why of those it could not read.
		status := a.accountsStatus(context.Background())
		a.mu.Lock()
		a.profileErrs = profileErrors{youtube: status.YouTube.Error, at: time.Now()}
		a.mu.Unlock()
		notify()
	}()
}

// accountsStatus is the state of the accounts with their profiles, which it
// asks the services for when it has not read them yet.
func (a *App) accountsStatus(ctx context.Context) AccountsStatus {
	return AccountsStatus{YouTube: a.youtubeAccount(ctx), Lastfm: a.lastfmAccount()}
}

// forgetSpotify deletes the Spotify sign-in MyTube 0.x kept, from the
// credential store and the files beside it, as MyTube plays from Spotify
// no more.
func forgetSpotify(dataDir string) {
	if err := newSecretStore(dataDir, spotifySecret, "spotify-token.json", "spotify-token.dpapi").Clear(); err != nil {
		log.Printf("mytube: forgetting the Spotify sign-in: %v", err)
	}
}
