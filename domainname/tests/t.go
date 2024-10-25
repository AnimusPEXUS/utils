package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/domainname"
)

func main() {
	for ii, i := range [][2]string{
		[2]string{"", "mozIlla.org"},
		[2]string{"mozilla.org", "mozIlla.org"},
		[2]string{"org", "mozilla.org"},
		[2]string{"mozilla.org", "org"},
		[2]string{"org", "com"},
		[2]string{"developer.mozilla.org", "mozilla.org"},
	} {

		first := domainname.NewDomainNameFromString(i[0])
		second := domainname.NewDomainNameFromString(i[1])

		fmt.Printf(`#%02d "%s" "%s"
         Is Equal?                   : %t,
         First Is Second's subdomain?: %t,
         Compare                     : %d,

`, ii, i[0], i[1],
			first.IsEqualTo(second),
			first.IsSubdomainTo(second),
			first.CompareTo(second),
		)
	}
}
