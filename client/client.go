package client

// Client represents a unique client attatched to a server
// The Aux can be used to store additional data
type Client struct {
	Active  bool
	UUID    uint32
	Aux     interface{}
	current string
	history []string
}

// SetBuffer updates the client's current buffer
// if set to "none", it marks it as inactive
// if buffer was previously inactive, it marks it as active
func (c *Client) SetBuffer(buffer string) {
	// Setting to inactive
	switch {
	case buffer == "none":
		if c.current != "none" {
			c.history = append(c.history, c.current)
		}

		c.Active = false
	// Coming out of inactive
	case c.current == "none":
		if buffer == "none" {
			c.Active = false
			return
		}

		c.Active = true
	case c.current == buffer:
		return
	default:
		c.Active = true
		c.history = append(c.history, c.current)
	}

	c.current = buffer
}

// Current returns the client's current buffer
func (c *Client) Current() string {
	return c.current
}

// History returns a list of all visited buffers
func (c *Client) History() []string {
	return c.history
}
