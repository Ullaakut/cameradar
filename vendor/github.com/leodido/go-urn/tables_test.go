package urn

import (
	"strconv"
	"strings"
)

func ierror(index int) string {
	return "Test case num. " + strconv.Itoa(index+1)
}

func herror(index int, test testCase) string {
	return ierror(index) + ", input \"" + string(test.in) + "\""
}

func rxpad(str string, lim int) string {
	str = str + strings.Repeat(" ", lim)
	return str[:lim]
}

type testCase struct {
	in   []byte // the input
	ok   bool   // whether it is valid or not
	obj  *URN   // a pointer to the resulting urn.URN instance
	str  string // string representation
	norm string // norm string representation
	estr string // error string
}

var tests = []testCase{
	// ok
	{
		[]byte("urn:simple:simple"),
		true,
		&URN{
			prefix: "urn",
			ID:     "simple",
			SS:     "simple",
		},
		"urn:simple:simple",
		"urn:simple:simple",
		"",
	},
	{
		[]byte("urn:ciao:%5D"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "%5D",
		},
		"urn:ciao:%5D",
		"urn:ciao:%5d",
		"",
	},

	// ok - RFC examples
	{
		[]byte("URN:foo:a123,456"),
		true,
		&URN{
			prefix: "URN",
			ID:     "foo",
			SS:     "a123,456",
		},
		"URN:foo:a123,456",
		"urn:foo:a123,456",
		"",
	},
	{
		[]byte("urn:foo:a123,456"),
		true,
		&URN{
			prefix: "urn",
			ID:     "foo",
			SS:     "a123,456",
		},
		"urn:foo:a123,456",
		"urn:foo:a123,456",
		"",
	},
	{
		[]byte("urn:FOO:a123,456"),
		true,
		&URN{
			prefix: "urn",
			ID:     "FOO",
			SS:     "a123,456",
		},
		"urn:FOO:a123,456",
		"urn:foo:a123,456",
		"",
	},
	{
		[]byte("urn:foo:A123,456"),
		true,
		&URN{
			prefix: "urn",
			ID:     "foo",
			SS:     "A123,456",
		},
		"urn:foo:A123,456",
		"urn:foo:A123,456",
		"",
	},
	{
		[]byte("urn:foo:a123%2C456"),
		true,
		&URN{
			prefix: "urn",
			ID:     "foo",
			SS:     "a123%2C456",
		},
		"urn:foo:a123%2C456",
		"urn:foo:a123%2c456",
		"",
	},
	{
		[]byte("URN:FOO:a123%2c456"),
		true,
		&URN{
			prefix: "URN",
			ID:     "FOO",
			SS:     "a123%2c456",
		},
		"URN:FOO:a123%2c456",
		"urn:foo:a123%2c456",
		"",
	},
	{
		[]byte("URN:FOO:ABC%FFabc123%2c456"),
		true,
		&URN{
			prefix: "URN",
			ID:     "FOO",
			SS:     "ABC%FFabc123%2c456",
		},
		"URN:FOO:ABC%FFabc123%2c456",
		"urn:foo:ABC%ffabc123%2c456",
		"",
	},
	{
		[]byte("URN:FOO:ABC%FFabc123%2C456%9A"),
		true,
		&URN{
			prefix: "URN",
			ID:     "FOO",
			SS:     "ABC%FFabc123%2C456%9A",
		},
		"URN:FOO:ABC%FFabc123%2C456%9A",
		"urn:foo:ABC%ffabc123%2c456%9a",
		"",
	},

	// ok - SCIM v2
	{
		[]byte("urn:ietf:params:scim:schemas:core:2.0:User"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ietf",
			SS:     "params:scim:schemas:core:2.0:User",
		},
		"urn:ietf:params:scim:schemas:core:2.0:User",
		"urn:ietf:params:scim:schemas:core:2.0:User",
		"",
	},
	{
		[]byte("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ietf",
			SS:     "params:scim:schemas:extension:enterprise:2.0:User",
		},
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		"",
	},
	{
		[]byte("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:userName"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ietf",
			SS:     "params:scim:schemas:extension:enterprise:2.0:User:userName",
		},
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:userName",
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:userName",
		"",
	},
	{
		[]byte("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:meta.lastModified"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ietf",
			SS:     "params:scim:schemas:extension:enterprise:2.0:User:meta.lastModified",
		},
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:meta.lastModified",
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:meta.lastModified",
		"",
	},

	// ok - minimum urn
	{
		[]byte("urn:a:b"),
		true,
		&URN{
			prefix: "urn",
			ID:     "a",
			SS:     "b",
		},
		"urn:a:b",
		"urn:a:b",
		"",
	},
	{
		[]byte("urn:a::"),
		true,
		&URN{
			prefix: "urn",
			ID:     "a",
			SS:     ":",
		},
		"urn:a::",
		"urn:a::",
		"",
	},
	{
		[]byte("urn:a:-"),
		true,
		&URN{
			prefix: "urn",
			ID:     "a",
			SS:     "-",
		},
		"urn:a:-",
		"urn:a:-",
		"",
	},

	// ok - URN prefix is case-insensitive
	{
		[]byte("URN:simple:simple"),
		true,
		&URN{
			prefix: "URN",
			ID:     "simple",
			SS:     "simple",
		},
		"URN:simple:simple",
		"urn:simple:simple",
		"",
	},
	{
		[]byte("Urn:simple:simple"),
		true,
		&URN{
			prefix: "Urn",
			ID:     "simple",
			SS:     "simple",
		},
		"Urn:simple:simple",
		"urn:simple:simple",
		"",
	},

	// ok - ID can contain the "urn" string but it can not be exactly equal to it
	{
		[]byte("urn:urna:simple"),
		true,
		&URN{
			prefix: "urn",
			ID:     "urna",
			SS:     "simple",
		},
		"urn:urna:simple",
		"urn:urna:simple",
		"",
	},
	{
		[]byte("urn:burnout:nss"),
		true,
		&URN{
			prefix: "urn",
			ID:     "burnout",
			SS:     "nss",
		},
		"urn:burnout:nss",
		"urn:burnout:nss",
		"",
	},
	{
		[]byte("urn:burn:nss"),
		true,
		&URN{
			prefix: "urn",
			ID:     "burn",
			SS:     "nss",
		},
		"urn:burn:nss",
		"urn:burn:nss",
		"",
	},
	{
		[]byte("urn:urnurnurn:x"),
		true,
		&URN{
			prefix: "urn",
			ID:     "urnurnurn",
			SS:     "x",
		},
		"urn:urnurnurn:x",
		"urn:urnurnurn:x",
		"",
	},

	// ok - ID can contains maximum 32 characters
	{
		[]byte("urn:abcdefghilmnopqrstuvzabcdefghilm:x"),
		true,
		&URN{
			prefix: "urn",
			ID:     "abcdefghilmnopqrstuvzabcdefghilm",
			SS:     "x",
		},
		"urn:abcdefghilmnopqrstuvzabcdefghilm:x",
		"urn:abcdefghilmnopqrstuvzabcdefghilm:x",
		"",
	},

	// ok - ID can be alpha numeric
	{
		[]byte("URN:123:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "123",
			SS:     "x",
		},
		"URN:123:x",
		"urn:123:x",
		"",
	},
	{
		[]byte("URN:1ab:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "1ab",
			SS:     "x",
		},
		"URN:1ab:x",
		"urn:1ab:x",
		"",
	},
	{
		[]byte("URN:a1b:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "a1b",
			SS:     "x",
		},
		"URN:a1b:x",
		"urn:a1b:x",
		"",
	},
	{
		[]byte("URN:a12:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "a12",
			SS:     "x",
		},
		"URN:a12:x",
		"urn:a12:x",
		"",
	},
	{
		[]byte("URN:cd2:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "cd2",
			SS:     "x",
		},
		"URN:cd2:x",
		"urn:cd2:x",
		"",
	},

	// ok - ID can contain an hyphen (not in its first position, see below)
	{
		[]byte("URN:abcd-:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "abcd-",
			SS:     "x",
		},
		"URN:abcd-:x",
		"urn:abcd-:x",
		"",
	},
	{
		[]byte("URN:abcd-abcd:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "abcd-abcd",
			SS:     "x",
		},
		"URN:abcd-abcd:x",
		"urn:abcd-abcd:x",
		"",
	},
	{
		[]byte("URN:a123-456z:x"),
		true,
		&URN{
			prefix: "URN",
			ID:     "a123-456z",
			SS:     "x",
		},
		"URN:a123-456z:x",
		"urn:a123-456z:x",
		"",
	},

	// ok - SS can contain the "urn" string, also be exactly equal to it
	{
		[]byte("urn:urnx:urn"),
		true,
		&URN{
			prefix: "urn",
			ID:     "urnx",
			SS:     "urn",
		},
		"urn:urnx:urn",
		"urn:urnx:urn",
		"",
	},
	{
		[]byte("urn:urnurnurn:urn"),
		true,
		&URN{
			prefix: "urn",
			ID:     "urnurnurn",
			SS:     "urn",
		},
		"urn:urnurnurn:urn",
		"urn:urnurnurn:urn",
		"",
	},
	{
		[]byte("urn:hey:urnurnurn"),
		true,
		&URN{
			prefix: "urn",
			ID:     "hey",
			SS:     "urnurnurn",
		},
		"urn:hey:urnurnurn",
		"urn:hey:urnurnurn",
		"",
	},

	// ok - SS can contains and discerns multiple colons, also at the end
	{
		[]byte("urn:ciao:a:b:c"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "a:b:c",
		},
		"urn:ciao:a:b:c",
		"urn:ciao:a:b:c",
		"",
	},
	{
		[]byte("urn:aaa:x:y:"),
		true,
		&URN{
			prefix: "urn",
			ID:     "aaa",
			SS:     "x:y:",
		},
		"urn:aaa:x:y:",
		"urn:aaa:x:y:",
		"",
	},
	{
		[]byte("urn:aaa:x:y:"),
		true,
		&URN{
			prefix: "urn",
			ID:     "aaa",
			SS:     "x:y:",
		},
		"urn:aaa:x:y:",
		"urn:aaa:x:y:",
		"",
	},

	// ok - SS can contain (and also start with) some non-alphabetical (ie., OTHER) characters
	{
		[]byte("urn:ciao:-"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "-",
		},
		"urn:ciao:-",
		"urn:ciao:-",
		"",
	},
	{
		[]byte("urn:ciao::"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     ":",
		},
		"urn:ciao::",
		"urn:ciao::",
		"",
	},
	{
		[]byte("urn:colon:::::nss"),
		true,
		&URN{
			prefix: "urn",
			ID:     "colon",
			SS:     "::::nss",
		},
		"urn:colon:::::nss",
		"urn:colon:::::nss",
		"",
	},
	{
		[]byte("urn:ciao:!"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "!",
		},
		"urn:ciao:!",
		"urn:ciao:!",
		"",
	},
	{
		[]byte("urn:ciao:!!*"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "!!*",
		},
		"urn:ciao:!!*",
		"urn:ciao:!!*",
		"",
	},
	{
		[]byte("urn:ciao:-!:-,:x"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "-!:-,:x",
		},
		"urn:ciao:-!:-,:x",
		"urn:ciao:-!:-,:x",
		"",
	},
	{
		[]byte("urn:ciao:=@"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "=@",
		},
		"urn:ciao:=@",
		"urn:ciao:=@",
		"",
	},
	{
		[]byte("urn:ciao:@!=%2C(xyz)+a,b.*@g=$_'"),
		true,
		&URN{
			prefix: "urn",
			ID:     "ciao",
			SS:     "@!=%2C(xyz)+a,b.*@g=$_'",
		},
		"urn:ciao:@!=%2C(xyz)+a,b.*@g=$_'",
		"urn:ciao:@!=%2c(xyz)+a,b.*@g=$_'",
		"",
	},

	// ok - SS can contain (and also start with) hexadecimal representation of octets
	{
		[]byte("URN:hexes:%25"),
		true,
		&URN{
			prefix: "URN",
			ID:     "hexes",
			SS:     "%25",
		},
		"URN:hexes:%25",
		"urn:hexes:%25",
		"",
	}, // Literal use of the "%" character in a namespace must be encoded using "%25"
	{
		[]byte("URN:x:abc%1Dz%2F%3az"),
		true,
		&URN{
			prefix: "URN",
			ID:     "x",
			SS:     "abc%1Dz%2F%3az",
		},
		"URN:x:abc%1Dz%2F%3az",
		"urn:x:abc%1dz%2f%3az",
		"",
	}, // Literal use of the "%" character in a namespace must be encoded using "%25"

	// no - ID can not start with an hyphen
	{
		[]byte("URN:-xxx:x"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},
	{
		[]byte("URN:---xxx:x"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},

	// no - ID can not start with a colon
	{
		[]byte("urn::colon:nss"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},
	{
		[]byte("urn::::nss"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},

	// no - ID can not contains more than 32 characters
	{
		[]byte("urn:abcdefghilmnopqrstuvzabcdefghilmn:specificstring"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 36]`,
	},

	// no - ID can not contain special characters
	{
		[]byte("URN:a!?:x"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 5]`,
	},
	{
		[]byte("URN:@,:x"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},
	{
		[]byte("URN:#,:x"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},
	{
		[]byte("URN:bc'.@:x"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 6]`,
	},

	// no - ID can not be equal to "urn"
	{
		[]byte("urn:urn:NSS"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to not contain the "urn" reserved string [col 7]`,
	},
	{
		[]byte("urn:URN:NSS"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to not contain the "urn" reserved string [col 7]`,
	},
	{
		[]byte("URN:URN:NSS"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to not contain the "urn" reserved string [col 7]`,
	},
	{
		[]byte("urn:UrN:NSS"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to not contain the "urn" reserved string [col 7]`,
	},
	{
		[]byte("urn:Urn:NSS"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to not contain the "urn" reserved string [col 7]`,
	},

	// no - ID can not contain spaces
	{
		[]byte("urn:white space:NSS"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 9]`,
	},

	// no - SS can not contain spaces
	{
		[]byte("urn:concat:no spaces"),
		false,
		nil,
		"",
		"",
		`expecting the specific string to be a string containing alnum, hex, or others ([()+,-.:=@;$_!*']) chars [col 13]`,
	},

	// no - SS can not contain reserved characters (can accept them only if %-escaped)
	{
		[]byte("urn:a:%"), // the presence of an "%" character in an URN MUST be followed by two characters from the <hex> character set
		false,
		nil,
		"",
		"",
		`expecting the specific string hex chars to be well-formed (%alnum{2}) [col 7]`,
	},
	{
		[]byte("urn:a:?"),
		false,
		nil,
		"",
		"",
		`expecting the specific string to be a string containing alnum, hex, or others ([()+,-.:=@;$_!*']) chars [col 6]`,
	},
	{
		[]byte("urn:a:#"),
		false,
		nil,
		"",
		"",
		`expecting the specific string to be a string containing alnum, hex, or others ([()+,-.:=@;$_!*']) chars [col 6]`,
	},
	{
		[]byte("urn:a:/"),
		false,
		nil,
		"",
		"",
		`expecting the specific string to be a string containing alnum, hex, or others ([()+,-.:=@;$_!*']) chars [col 6]`,
	},

	// no - Incomplete URNs
	{
		[]byte("urn:"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},
	{
		[]byte("urn::"),
		false,
		nil,
		"",
		"",
		`expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its start) [col 4]`,
	},
	{
		[]byte("urn:a:"),
		false,
		nil,
		"",
		"",
		`expecting the specific string to be a string containing alnum, hex, or others ([()+,-.:=@;$_!*']) chars [col 6]`,
	},
	// {
	// 	"urn:a",
	// 	false,
	// 	nil,
	// 	"",
	// 	"",
	// 	"",
	// },
}

var equivalenceTests = []struct {
	eq bool
	lx []byte
	rx []byte
}{
	{
		true,
		[]byte("urn:foo:a123%2C456"),
		[]byte("URN:FOO:a123%2c456"),
	},
	{
		true,
		[]byte("urn:foo:AbC123%2C456"),
		[]byte("URN:FOO:AbC123%2c456"),
	},
	{
		true,
		[]byte("urn:foo:AbC123%2C456%1f"),
		[]byte("URN:FOO:AbC123%2c456%1f"),
	},
	{
		true,
		[]byte("URN:foo:a123,456"),
		[]byte("urn:foo:a123,456"),
	},
	{
		true,
		[]byte("URN:foo:a123,456"),
		[]byte("urn:FOO:a123,456"),
	},
	{
		true,
		[]byte("urn:foo:a123,456"),
		[]byte("urn:FOO:a123,456"),
	},
	{
		true,
		[]byte("urn:ciao:%2E"),
		[]byte("urn:ciao:%2e"),
	},
	{
		false,
		[]byte("urn:foo:A123,456"),
		[]byte("URN:foo:a123,456"),
	},
	{
		false,
		[]byte("urn:foo:A123,456"),
		[]byte("urn:foo:a123,456"),
	},
	{
		false,
		[]byte("urn:foo:A123,456"),
		[]byte("urn:FOO:a123,456"),
	},
}
