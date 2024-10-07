package test

import (
	"log"
	"testing"
)

func SetupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

// Almost the same as the above, but this one is for single test instead of collection of tests
func SetupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
	}
}
