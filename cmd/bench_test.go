// Copyright © 2017 Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"github.com/dougEfresh/passwd-pot/api"
	"testing"
)

func BenchmarkEvent(b *testing.B) {
	var event api.Event
	b.ReportAllocs()
	if err := json.Unmarshal([]byte(requestBodyOrigin), &event); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkLookup(b *testing.B) {
	var event api.Event
	b.ReportAllocs()
	if err := json.Unmarshal([]byte(requestBodyOrigin), &event); err != nil {
		b.Fatal(err)
	}
	//id, _ := defaultEventClient.recordEvent(event)
	//event.ID = id
	for i := 0; i < b.N; i++ {
		//	defaultEventClient.resolveGeoEvent(event)
	}

}

func BenchmarkGeoCache(b *testing.B) {
	b.ReportAllocs()
	geoCache.set("blah", 1)
	for i := 0; i < b.N; i++ {
		geoCache.get("blah")
	}
}
