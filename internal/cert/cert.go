package cert

import _ "embed"

//go:embed secret.pem
var RawPrivKey []byte
