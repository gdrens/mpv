package mpv

import (
	"errors"
	"fmt"
)

// Client is a more comfortable higher level interface
// to LLClient. It can use any LLClient implementation.
type Client struct {
	LLClient
}

// NewClient creates a new highlevel client based on a lowlevel client.
func NewClient(llClient LLClient) *Client {
	return &Client{
		LLClient: llClient,
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
func (c *Client) LoadFile(path string, mode string) error {
	if mode == "" {
		mode = "append-play"
	}
	_, err := c.Exec("loadfile", path, mode)
	return err
}

// Mode options for Seek
const (
	SeekModeRelative = "relative"
	SeekModeAbsolute = "absolute"
)

// Seek seeks to a position in the current file.
func (c *Client) Seek(n int) error {
	_, err := c.Exec("seek", n)
	return err
}

// PlaylistNext plays the next playlistitem or NOP if no item is available.
func (c *Client) PlayNext() error {
	_, err := c.Exec("playlist-next")
	return err
}

// PlaylistPrevious plays the previous playlistitem or NOP if no item is available.
func (c *Client) PlayPrev() error {
	_, err := c.Exec("playlist-prev")
	return err
}

// Return Playlist Current Pos
func (c *Client) PlayPos() int {
	n, _ := c.GetFloatProperty("playlist-pos")
	return int(n)
}

// Return Playlist
func (c *Client) Playlist() []string {
	var names []string
	resp, _ := c.Exec("get_property", "playlist")
	for _, v := range resp.Data.([]interface{}) {
		for k, v1 := range v.(map[string]interface{}) {
			if k == "filename" {
				names = append(names, v1.(string))
			}
		}
	}
	return names
}

// Remove current Playlist
func (c *Client) PlayRemove() error {
	_, err := c.Exec("playlist-remove")
	return err
}

// Remove the specified playlistitem
func (c *Client) PlayIndexRemove(n int) error {
	_, err := c.Exec("playlist-remove", n)
	return err
}

// Clear Playlist (keep the playing)
func (c *Client) PlayClear() error {
	_, err := c.Exec("playlist-clear")
	return err
}

// Play the specified item
func (c *Client) PlayIndex(n int) error {
	_, err := c.Exec("playlist-play-index", n)
	return err
}

// loop-playlist
func (c *Client) PlayLoop() error {
	return c.SetProperty("loop-playlist", true)
}

// Unloop-playlist
func (c *Client) PlayUnLoop() error {
	return c.SetProperty("loop-playlist", false)
}

// Return playloop status
func (c *Client) IsPlayLoop() bool {
	b := c.GetProperty("loop-playlist")
	return b == "inf"
}

// Shuffle the playlist
func (c *Client) PlayShuffle() error {
	_, err := c.Exec("playlist-suffle")
	return err
}

// UnShuffle the playlist
func (c *Client) PlayUnShuffle() error {
	_, err := c.Exec("playlist-unshuffle")
	return err
}

// Return the playlist-count
func (c *Client) PlaylistCount() int {
	n, _ := c.GetFloatProperty("playlist-count")
	return int(n)
}

// Mode options for LoadList
const (
	LoadListModeReplace = "replace"
	LoadListModeAppend  = "append"
)

// LoadList loads a playlist from path. It can either replace the current playlist (LoadListModeReplace)
// or append to the current playlist (LoadListModeAppend).
func (c *Client) LoadList(path string, mode string) error {
	if mode == "" {
		mode = "replace"
	}
	_, err := c.Exec("loadlist", path, mode)
	return err
}

// GetProperty reads a property by name and returns the data as a string.
func (c *Client) GetProperty(name string) string {
	res, _ := c.Exec("get_property", name)
	if res == nil {
		return ""
	}
	if v, ok := res.Data.(string); ok {
		return v
	}
	return fmt.Sprintf("%#v", res.Data)
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
func (c *Client) CurrentFile() string {
	return c.GetProperty("filename")
}

// Path returns the currently playing path
func (c *Client) CurerentFileWithPath() string {
	return c.GetProperty("path")
}

// Pause returns true if the player is paused
func (c *Client) IsPause() bool {
	b, _ := c.GetBoolProperty("pause")
	return b
}

// SetPause pauses or unpauses the player
func (c *Client) Pause() error {
	_, err := c.Exec("cycle", "pause")
	return err
}

// Idle returns true if the player is idle
func (c *Client) IsIdle() bool {
	b, _ := c.GetBoolProperty("idle")
	return b
}

// Mute returns true if the player is muted.
func (c *Client) IsMute() bool {
	b, _ := c.GetBoolProperty("mute")
	return b
}

// SetMute mutes or unmutes the player.
func (c *Client) Mute() error {
	_, err := c.Exec("cycle", "mute")
	return err
}

// Fullscreen returns true if the player is in fullscreen mode.
func (c *Client) IsFullscreen() bool {
	b, _ := c.GetBoolProperty("fullscreen")
	return b
}

// SetFullscreen activates/deactivates the fullscreen mode.
func (c *Client) Fullscreen() error {
	_, err := c.Exec("cycle", "fullscreen")
	return err
}

// Volume returns the current volume level.
func (c *Client) CurrentVolume() int {
	v, _ := c.GetFloatProperty("volume")
	return int(v)
}

// Set Volume level
func (c *Client) Volume(level int) error {
	return c.SetProperty("volume", level)
}

// Speed returns the current playback speed.
func (c *Client) CurrentSpeed() float64 {
	s, _ := c.GetFloatProperty("speed")
	return s
}

// Set playback speed
func (c *Client) Speed(n float64) error {
	return c.SetProperty("speed", n)
}

// Duration returns the duration of the currently playing file.
func (c *Client) Duration() float64 {
	v, _ := c.GetFloatProperty("duration")
	return v
}

// Position returns the current playback position in seconds.
func (c *Client) Position() float64 {
	b, _ := c.GetFloatProperty("time-pos")
	return b
}

// PercentPosition returns the current playback position in percent.
func (c *Client) PercentPosition() float64 {
	b, _ := c.GetFloatProperty("percent-pos")
	return b
}

// Register Event HandFunc
func (c *Client) RegisterEvent(eventName string, handle func()) {
	c.LLClient.RegisterEvent(eventName, handle)
}

// loop-file
func (c *Client) FileLoop() error { //"inf" is Infinite loop
	return c.SetProperty("loop-file", true)
}

// Unloop-file
func (c *Client) FileUnLoop() error { //"inf" is Infinite loop
	return c.SetProperty("loop-file", false)
}

// Return fileloop status
func (c *Client) IsFileLoop() bool {
	b := c.GetProperty("loop-file")
	return b == "inf"
}

// time-remaining
func (c *Client) TimeRemaining() float64 {
	b, _ := c.GetFloatProperty("time-remaining")
	return b
}

// Playlist shuffle
func (c *Client) Shuffle() error {
	_, err := c.Exec("cycle", "shuffle")
	return err
}

// Return shuffle status
func (c *Client) IsShuffle() bool {
	b, _ := c.GetBoolProperty("shuffle")
	return b
}

// Get file-format
func (c *Client) Format() string {
	return c.GetProperty("file-format")
}

// Get Audio-bitrate
func (c *Client) AudioBitrate() int {
	v, _ := c.GetFloatProperty("audio-bitrate")
	return int(v)
}

// Get Video-bitrate
func (c *Client) VideoBitrate() int {
	v, _ := c.GetFloatProperty("video-bitrate")
	return int(v)
}

// Get file-size
func (c *Client) FileSize() int {
	v, _ := c.GetFloatProperty("file-size")
	return int(v)
}

// Get media-title
func (c *Client) MediaTitle() string {
	return c.GetProperty("media-title")
}

// Quit
func (c *Client) Quit() error {
	_, err := c.Exec("quit")
	return err
}

// Stop and clear playlist
func (c *Client) Stop() error {
	_, err := c.Exec("stop")
	return err
}
