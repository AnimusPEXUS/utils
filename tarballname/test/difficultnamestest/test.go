package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tarballname"
)

var DIFFICULT_NAMES []string

func init() {
	DIFFICULT_NAMES = []string{
		"GeSHi-1.0.2-beta-1.tar.bz2",
		"Perl-Dist-Strawberry-BuildPerl-5101-2.11_10.tar.gz",
		"bind-9.9.1-P2.tar.gz",
		"boost_1_25_1.tar.bz2",
		"dahdi-linux-complete-2.1.0.3+2.1.0.2.tar.gz",
		"dhcp-4.1.2rc1.tar.gz",
		"dvd+rw-tools-5.5.4.3.4.tar.gz",
		"fontforge_full-20120731-b.tar.bz2",
		"gtk+-3.12.0.tar.xz",
		"lynx2.8.7rel.1.tar.bz2",
		"name.tar.gz",
		"ogre_src_v1-8-1.tar.bz2",
		"openjdk-8-src-b132-03_mar_2014.zip",
		"openssl-0.9.7a.tar.gz",
		"org.apache.felix.ipojo.manipulator-1.8.4-project.tar.gz",
		"pa_stable_v19_20140130.tgz",
		"pkcs11-helper-1.05.tar.bz2",
		"qca-pkcs11-0.1-20070425.tar.bz2",
		"tcl8.4.19-src.tar.gz",
		"wmirq-0.1-source.tar.gz",
		"xc-1.tar.gz",
		"xf86-input-acecad-1.5.0.tar.bz2",
		"xf86-input-elo2300-1.1.2.tar.bz2",
		"ziplock-1.7.3-source-release.zip",
		// delimiters missing between version numbers :-
		"unzip60.tar.gz",
		"zip30.tar.gz",
		"zip30c.tar.gz",
	}
}

func main() {

	for ii, i := range DIFFICULT_NAMES {

		fmt.Println("----------------------")

		fmt.Printf("test: #%02d, \"%s\"\n", ii, i)

		res, err := tarballname.Parse(i)

		if err != nil {
			fmt.Printf("exception: %s\n", err)
			continue
		}

		if res == nil {
			fmt.Println("couldn't parse")
			continue
		}

		fmt.Printf("parser returned:\n")
		fmt.Printf("%s\n", res)

	}

}
