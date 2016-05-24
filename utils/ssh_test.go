package utils_test

import (
	"fmt"
	"github.com/evq/chunnel/utils"
)

func ExampleParseUserHost_all() {
	fmt.Println(utils.ParseUserHost("root@test:8000"))
	// Output: root test 8000 <nil>
}

func ExampleParseUserHost_invalid() {
	fmt.Println(utils.ParseUserHost("root@test:n"))
	// Output: root  0 strconv.ParseInt: parsing "n": invalid syntax
}

func ExampleParseUserHost_noUser() {
	fmt.Println(utils.ParseUserHost("test:8000"))
	// Output: ev test 8000 <nil>
}

func ExampleParseUserHost_noPort() {
	fmt.Println(utils.ParseUserHost("root@test"))
	// Output: root test 22 <nil>
}

func ExampleParseUserHost_fqdn() {
	fmt.Println(utils.ParseUserHost("foo.pvt.test.com"))
	// Output: ev foo.pvt.test.com 22 <nil>
}
