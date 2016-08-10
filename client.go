package mpv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Client is a more comfortable higher level interface
// to LLClient. It can use any LLClient implementation.
type Client struct {
	LLClient
}

// NewClient creates a new highlevel client based on a lowlevel client.
func NewClient(llclient LLClient) *Client {
	return &Client{
		llclient,
	}
}

// Mode options for Loadfile
const (
	LoadFileModeReplace    = "replace"
	LoadFileModeAppend     = "append"
	LoadFileModeAppendPlay = "append-play" // Starts if nothing is playing
)

// Loadfile loads a file, it either replaces the currently playing file (LoadFileModeReplace),
// appends to the current playlist (LoadFileModeAppend) or appends to playlist and plays if
// nothing is playing right now (LoadFileModeAppendPlay)
func (c *Client) Loadfile(path string, mode string) error {
	_, err := c.Exec("loadfile", path, mode)
	return err
}

// Mode options for Seek
const (
	SeekModeRelative = "relative"
	SeekModeAbsolute = "absolute"
)

//mpv PlayList struct
type Playlist struct {
	FileName string `json:"filename"`
	Current  bool   `json:"current"`
	Playing  bool   `json:"playing"`
	Id       int    `json:"id"`
}

// Seek seeks to a position in the current file.
// Use mode to seek relative to current position (SeekModeRelative) or absolute (SeekModeAbsolute).
func (c *Client) Seek(n int, mode string) error {
	return c.InputCmd("seek", strconv.Itoa(n), mode)
}

// PlaylistNext plays the next playlistitem or NOP if no item is available.
func (c *Client) PlaylistNext() error {
	return c.InputCmd("playlist-next", "weak")
}

// PlaylistPrevious plays the previous playlistitem or NOP if no item is available.
func (c *Client) PlaylistPrevious() error {
	return c.InputCmd("playlist-prev", "weak")
}

//Get Playlist Current Pos
func (c *Client) PlaylistCurrentPos() (float64, error) {
	return c.GetFloatProperty("playlist-current-pos")
}

//Get Playlist
func (c *Client) Playlist() []Playlist {
	tmp := new(struct {
		Data []Playlist
	})
	resp, err := c.GetRawProperty("playlist")
	if err != nil {
		return nil
	}
	err = json.Unmarshal(resp, tmp)
	if err != nil {
		return nil
	}
	return tmp.Data
}

//Remove Playlist
func (c *Client) PlaylistRemove(n int) error {
	return c.InputCmd("raw", fmt.Sprintf("\"command\":[\"playlist-remove\"],\"index\":%d", n))
}

//Clear Playlist
func (c *Client) PlaylistClear() error {
	return c.InputCmd("playlist-clear")
}

//Playlist Play index
func (c *Client) PlaylistPlayIndex(n int) error {
	return c.InputCmd("playlist-play-index", n)
}

//playlist-shuffle
func (c *Client) PlaylistShuffle() error {
	return c.InputCmd("playlist-shuffle")
}

//playlist-unshuffle
func (c *Client) PlaylistUnShuffle() error {
	return c.InputCmd("playlist-unshuffle")
}

//loop-playlist
func (c *Client) PlaylistLoop(n string) error { //"inf" is Infinite loop
	return c.InputCmd("set", "loop-playlist", n)
}

//playlist-count
func (c *Client) PlaylistCount() (float64, error) {
	return c.GetFloatProperty("playlist-count")
}

// Mode options for LoadList
const (
	LoadListModeReplace = "replace"
	LoadListModeAppend  = "append"
)

// LoadList loads a playlist from path. It can either replace the current playlist (LoadListModeReplace)
// or append to the current playlist (LoadListModeAppend).
func (c *Client) LoadList(path string, mode string) error {
	return c.InputCmd("loadlist", path, mode)
}

//GetProperty reads a property by name and returns the data as Raw
func (c *Client) GetRawProperty(name string) ([]byte, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return nil, err
	}
	return res.Bytes, err

}

// GetProperty reads a property by name and returns the data as a string.
func (c *Client) GetProperty(name string) (string, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return "", err
	}
	return fmt.Sprintf("%#v", res.Data), err
}

// SetProperty sets the value of a property.
func (c *Client) SetProperty(name string, value interface{}) error {
	_, err := c.Exec("set_property", name, value)
	return err
}

// ErrInvalidType is returned if the response data does not match the methods return type.
// Use GetProperty or find matching type in mpv docs.
var ErrInvalidType = errors.New("Invalid type")

// GetFloatProperty reads a float property and returns the data as a float64.
func (c *Client) GetFloatProperty(name string) (float64, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return 0, err
	}
	if val, found := res.Data.(float64); found {
		return val, err
	}
	return 0, ErrInvalidType
}

// GetBoolProperty reads a bool property and returns the data as a boolean.
func (c *Client) GetBoolProperty(name string) (bool, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return false, err
	}
	if val, found := res.Data.(bool); found {
		return val, err
	}
	return false, ErrInvalidType
}

// Filename returns the currently playing filename
func (c *Client) Filename() (string, error) {
	return c.GetProperty("filename")
}

// Path returns the currently playing path
func (c *Client) Path() (string, error) {
	return c.GetProperty("path")
}

// Pause returns true if the player is paused
func (c *Client) Pause() (bool, error) {
	return c.GetBoolProperty("pause")
}

// SetPause pauses or unpauses the player
func (c *Client) SetPause() error {
	return c.InputCmd("cycle", "pause")
}

// Idle returns true if the player is idle
func (c *Client) Idle() (bool, error) {
	return c.GetBoolProperty("idle")
}

// Mute returns true if the player is muted.
func (c *Client) Mute() (bool, error) {
	return c.GetBoolProperty("mute")
}

// SetMute mutes or unmutes the player.
func (c *Client) SetMute(mute bool) error {
	return c.SetProperty("mute", mute)
}

// Fullscreen returns true if the player is in fullscreen mode.
func (c *Client) Fullscreen() (bool, error) {
	return c.GetBoolProperty("fullscreen")
}

// SetFullscreen activates/deactivates the fullscreen mode.
func (c *Client) SetFullscreen(v bool) error {
	return c.SetProperty("fullscreen", v)
}

// Volume returns the current volume level.
func (c *Client) Volume() (float64, error) {
	return c.GetFloatProperty("volume")
}

//Set Volume level
func (c *Client) SetVolume(level int) error {
	return c.SetProperty("volume", level)
}

// Speed returns the current playback speed.
func (c *Client) Speed() (float64, error) {
	return c.GetFloatProperty("speed")
}

//Set Speed
func (c *Client) SetSpeed(n int) error {
	return c.InputCmd("set", "speed", fmt.Sprintf("%d", n))
}

// Duration returns the duration of the currently playing file.
func (c *Client) Duration() (float64, error) {
	return c.GetFloatProperty("duration")
}

// Position returns the current playback position in seconds.
func (c *Client) Position() (float64, error) {
	return c.GetFloatProperty("time-pos")
}

// PercentPosition returns the current playback position in percent.
func (c *Client) PercentPosition() (float64, error) {
	return c.GetFloatProperty("percent-pos")
}

//Register Event HandFunc
func (c *Client) RegisterEvent(eventName string, handle handleEvent) {
	c.LLClient.(*IPCClient).registerEvent(eventName, handle)
}

//loop-file
func (c *Client) FileLoop(n string) error { //"inf" is Infinite loop
	return c.InputCmd("set", "loop-file", n)
}

//time-remaining
func (c *Client) TimeRemaining() (float64, error) {
	return c.GetFloatProperty("time-remaining")
}

//shuffle
func (c *Client) Shuffle() error {
	return c.InputCmd("cycle", "shuffle")
}

//Get shuffle status
func (c *Client) GetShuffle() (bool, error) {
	return c.GetBoolProperty("shuffle")
}

//Quit
func (c *Client) Quit() error {
	return c.InputCmd("quit")
}

//Stop and clear playlist
func (c *Client) Stop() error {
	return c.InputCmd("stop")
}

//Input Cmd
func (c *Client) InputCmd(cmd ...interface{}) error {
	resp, err := c.Exec(cmd...)
	if resp != nil {
		return errors.New(resp.Err)
	}
	return err
}
