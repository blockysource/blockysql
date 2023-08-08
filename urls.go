// Copyright 2023 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blockysql

import (
	"context"
	"net/url"

	"github.com/blockysource/go-pkg/urlopener"
)

var defaultURLMux = new(URLMux)

// DefaultURLMux returns the URLMux used by OpenSender.
//
// Driver packages can use this to register their SenderURLOpener on the mux.
func DefaultURLMux() *URLMux {
	return defaultURLMux
}

// DBURLOpener is an interface that allows to open a database connection
// for a given URL.
type DBURLOpener interface {
	// OpenDBURL opens a database connection for the given URL.
	OpenDBURL(ctx context.Context, u *url.URL) (*DB, error)
}

// URLMux is URL opener multiplexer. It matches the scheme of the URLs against
// a set of registered schemes and calls the opener that matches the URL's
// scheme. See https://gocloud.dev/concepts/urls/ for more information.
type URLMux struct {
	dbSchemes urlopener.SchemeMap
}

// DBSchemes returns a sorted slice of the registered DB schemes.
func (mux *URLMux) DBSchemes() []string { return mux.dbSchemes.Schemes() }

// ValidDBScheme returns true iff scheme has been registered for MailProviders.
func (mux *URLMux) ValidDBScheme(scheme string) bool {
	return mux.dbSchemes.ValidScheme(scheme)
}

// RegisterDB registers the opener with the given scheme. If an opener
// already exists for the scheme, RegisterDB panics.
func (mux *URLMux) RegisterDB(scheme string, opener DBURLOpener) {
	mux.dbSchemes.Register("emails", "DB", scheme, opener)
}

// OpenDB calls OpenDBURL with the URL parsed from urlstr.
// OpenDB is safe to call from multiple goroutines.
func (mux *URLMux) OpenDB(ctx context.Context, urlstr string) (*DB, error) {
	opener, u, err := mux.dbSchemes.FromString("DB", urlstr)
	if err != nil {
		return nil, err
	}
	return opener.(DBURLOpener).OpenDBURL(ctx, u)
}

// OpenDB opens the database identified by the URL given.
// i.e. "mysql://user:password@localhost:3306/database"
// OpenDB is safe to call from multiple goroutines.
func OpenDB(ctx context.Context, urlstr string) (*DB, error) {
	return defaultURLMux.OpenDB(ctx, urlstr)
}
