package db

type LoginType int

// Note: new type must append to the end of list to maintain compatibility.
const (
	LoginNotype LoginType = iota
	LoginPlain            // 1
	LoginLdap             // 2
	LoginSmtp             // 3
	LoginPam              // 4
	LoginDldap            // 5
	LoginGithub           // 6
	LoginWeixin			//7
	LoginPhone			//8
)