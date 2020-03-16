package domainname

import (
	"strings"
)

type DomainName struct {
	parts        []string
	recalculated string
}

func NewDomainNameFromStringsSlice(domainname []string) *DomainName {
	self := &DomainName{}
	self.SetFromStringsSlice(domainname)
	return self
}

func NewDomainNameFromString(domainname string) *DomainName {
	self := &DomainName{}
	self.SetFromString(domainname)
	return self
}

func (self *DomainName) SetFromString(domainname string) {
	t := strings.Split(domainname, ".")
	self.SetFromStringsSlice(t)
	return
}

func (self *DomainName) SetFromStringsSlice(domainname []string) {
	self.parts = domainname
	for i := 0; i != len(self.parts); i++ {
		self.parts[i] = strings.ToLower(self.parts[i])
	}
	self.recalculated = strings.Join(self.parts, ".")
	return
}

func (self *DomainName) IsEqualTo(to *DomainName) bool {

	len_self := len(self.parts)
	len_to := len(to.parts)

	if len_self != len_to {
		return false
	}

	for i := 0; i != len_self; i++ {
		if self.parts[i] != to.parts[i] {
			return false
		}
	}

	return true
}

// Check is self subdomain to two
func (self *DomainName) IsSubdomainTo(to *DomainName) bool {

	len_self := len(self.parts)
	len_to := len(to.parts)

	if len_self <= len_to {
		return false
	}

	if !NewDomainNameFromStringsSlice(self.parts[len_self-len_to:]).IsEqualTo(to) {
		return false
	}

	return true
}

func (self *DomainName) CompareTo(to *DomainName) int {

	len_self := len(self.parts)
	len_to := len(to.parts)

	len_max := len_self
	if len_to > len_max {
		len_max = len_to
	}

	for i := 0; i != len_max; i++ {

		i_self := len_self - i - 1
		i_to := len_to - i - 1

		if i_self < 0 || i_to < 0 {
			if len_self > len_to {
				return 1
			}
			if len_self < len_to {
				return -1
			}
		}

		if self.parts[i_self] > to.parts[i_to] {
			return 1
		}
		if self.parts[i_self] < to.parts[i_to] {
			return -1
		}
	}

	return 0
}

func (self *DomainName) String() string {
	return self.recalculated
}

func (self *DomainName) Len() int {
	return len(self.parts)
}

func (self *DomainName) Item(i int) string {
	return self.parts[i]
}

func (self *DomainName) IsEqualToString(to string) bool {
	t := NewDomainNameFromString(to)
	return self.IsEqualTo(t)
}

func (self *DomainName) IsSubdomainToString(to string) bool {
	t := NewDomainNameFromString(to)
	return self.IsSubdomainTo(t)
}

func (self *DomainName) CompareToString(to string) int {
	t := NewDomainNameFromString(to)
	return self.CompareTo(t)
}

func IsEqualTo(one, two string) bool {
	// NOTE: this is intentioanlly not a simple strings comparison: DomainName
	//       constructors may gain additional checks and conversions in the
	//       future
	o := NewDomainNameFromString(one)
	t := NewDomainNameFromString(two)
	return o.IsEqualTo(t)
}

// Check is one subdomain to two
func IsSubdomainTo(one, two string) bool {
	o := NewDomainNameFromString(one)
	t := NewDomainNameFromString(two)
	return o.IsSubdomainTo(t)
}

func Compare(one, two string) int {
	o := NewDomainNameFromString(one)
	t := NewDomainNameFromString(two)
	return o.CompareTo(t)
}
