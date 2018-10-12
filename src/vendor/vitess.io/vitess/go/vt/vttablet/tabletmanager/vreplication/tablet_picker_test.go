/*
Copyright 2018 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vreplication

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

func TestPickSimple(t *testing.T) {
	ts := createTopo()
	defer ts.Close()
	want := addTablet(ts, 100, "0", topodatapb.TabletType_REPLICA, true, true)

	tp, err := newTabletPicker(ts, testCell, testKeyspace, testShard, "replica")
	if err != nil {
		t.Fatal(err)
	}
	defer tp.Close()

	tablet, err := tp.Pick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(want, tablet) {
		t.Errorf("Pick: %v, want %v", tablet, want)
	}
}

func TestPickFromTwoHealthy(t *testing.T) {
	ts := createTopo()
	defer ts.Close()
	want1 := addTablet(ts, 100, "0", topodatapb.TabletType_REPLICA, true, true)
	want2 := addTablet(ts, 101, "0", topodatapb.TabletType_RDONLY, true, true)

	tp, err := newTabletPicker(ts, testCell, testKeyspace, testShard, "replica,rdonly")
	if err != nil {
		t.Fatal(err)
	}
	defer tp.Close()

	tablet, err := tp.Pick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(tablet, want1) {
		t.Errorf("Pick:\n%v, want\n%v", tablet, want1)
	}

	tp, err = newTabletPicker(ts, testCell, testKeyspace, testShard, "rdonly,replica")
	if err != nil {
		t.Fatal(err)
	}
	defer tp.Close()

	tablet, err = tp.Pick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(tablet, want2) {
		t.Errorf("Pick:\n%v, want\n%v", tablet, want2)
	}
}

func TestPickFromSomeUnhealthy(t *testing.T) {
	ts := createTopo()
	defer ts.Close()
	_ = addTablet(ts, 100, "0", topodatapb.TabletType_REPLICA, false, false)
	want := addTablet(ts, 101, "0", topodatapb.TabletType_RDONLY, false, true)

	tp, err := newTabletPicker(ts, testCell, testKeyspace, testShard, "replica,rdonly")
	if err != nil {
		t.Fatal(err)
	}
	defer tp.Close()

	tablet, err := tp.Pick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(tablet, want) {
		t.Errorf("Pick:\n%v, want\n%v", tablet, want)
	}
}

func TestPickError(t *testing.T) {
	ts := createTopo()
	defer ts.Close()
	_ = addTablet(ts, 100, "0", topodatapb.TabletType_REPLICA, false, false)

	_, err := newTabletPicker(ts, testCell, testKeyspace, testShard, "badtype")
	want := "failed to parse list of tablet types: badtype"
	if err == nil || err.Error() != want {
		t.Errorf("newTabletPicker err: %v, want %v", err, want)
	}

	tp, err := newTabletPicker(ts, testCell, testKeyspace, testShard, "replica,rdonly")
	if err != nil {
		t.Fatal(err)
	}
	defer tp.Close()

	_, err = tp.Pick(context.Background())
	want = "can't find any healthy source tablet for ks 0 [REPLICA RDONLY]"
	if err == nil || err.Error() != want {
		t.Errorf("Pick err: %v, want %v", err, want)
	}
}
