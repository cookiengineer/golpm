package lpm

import "testing"

func assertHashMapSearch(t *testing.T, hashmap HashMap, value string, want_prefix uint8, want_address string) {

	t.Helper()

	got := hashmap.Search(value)

	if got == nil {

		t.Errorf("Search(%q) returned nil, expected %s/%d", value, want_address, want_prefix)
		return

	}

	if got.Prefix != want_prefix {
		t.Errorf("Search(%q).Prefix = %d, expected %d", value, got.Prefix, want_prefix)
	}

	if got.Address != want_address {
		t.Errorf("Search(%q).Address = %q, expected %q", value, got.Address, want_address)
	}

}

func TestHashMapInsertWrongOrderIPv4(t *testing.T) {

	hashmap := NewHashMap()

	if hashmap.Insert("1.2.3.0/24") != true {
		t.Errorf("Insert(1.2.3.0/24) returned false, expected true")
	}

	if hashmap.Insert("1.2.3.0/16") != true {
		t.Errorf("Insert(1.2.3.0/16) returned false, expected true")
	}

	if hashmap.Insert("1.2.3.0/8") != true {
		t.Errorf("Insert(1.2.3.0/8) returned false, expected true")
	}

	assertHashMapSearch(t, hashmap, "1.2.3.4", 24, "1.2.3.0")

}

func TestHashMapInsertCorrectOrderIPv4(t *testing.T) {

	hashmap := NewHashMap()

	hashmap.Insert("1.2.3.0/8")
	hashmap.Insert("1.2.3.0/16")
	hashmap.Insert("1.2.3.0/24")

	assertHashMapSearch(t, hashmap, "1.2.3.4", 24, "1.2.3.0")

}

func TestHashMapSearchLongestPrefixMatch(t *testing.T) {

	hashmap := NewHashMap()

	hashmap.Insert("192.168.0.0/24")
	hashmap.Insert("192.169.128.0/17")
	hashmap.Insert("192.169.0.0/16")
	hashmap.Insert("192.169.0.0/24")

	assertHashMapSearch(t, hashmap, "192.168.0.123", 24, "192.168.0.0")
	assertHashMapSearch(t, hashmap, "192.169.0.123", 24, "192.169.0.0")
	assertHashMapSearch(t, hashmap, "192.169.128.123", 17, "192.169.128.0")
	assertHashMapSearch(t, hashmap, "192.169.64.123", 16, "192.169.0.0")

	got := hashmap.Search("8.8.8.8")

	if got != nil {
		t.Errorf("Search(8.8.8.8) = %s/%d, expected nil", got.Address, got.Prefix)
	}

}

func TestHashMapDuplicateInsert(t *testing.T) {

	hashmap := NewHashMap()

	if hashmap.Insert("1.2.3.0/24") != true {
		t.Errorf("first Insert(1.2.3.0/24) returned false, expected true")
	}

	if hashmap.Insert("1.2.3.0/24") != false {
		t.Errorf("duplicate Insert(1.2.3.0/24) returned true, expected false")
	}

}

func TestHashMapSamePrefixDifferentAddress(t *testing.T) {

	hashmap := NewHashMap()

	hashmap.Insert("1.2.3.0/24")
	hashmap.Insert("1.2.4.0/24")

	assertHashMapSearch(t, hashmap, "1.2.3.4", 24, "1.2.3.0")
	assertHashMapSearch(t, hashmap, "1.2.4.4", 24, "1.2.4.0")

}

func TestHashMapInsertWrongOrderIPv6(t *testing.T) {

	hashmap := NewHashMap()

	hashmap.Insert("2001:0db8:0000:0000:0000:0000:0000:0000/64")
	hashmap.Insert("2001:0db8:0000:0000:0000:0000:0000:0000/48")
	hashmap.Insert("2001:0db8:0000:0000:0000:0000:0000:0000/32")

	assertHashMapSearch(t, hashmap, "2001:0db8:0000:0000:0000:0000:0000:0001", 64, "[2001:0db8:0000:0000:0000:0000:0000:0000]")

}
