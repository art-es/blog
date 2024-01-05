package testutil

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

type uuidMatcher struct{}

func (m *uuidMatcher) Matches(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return false
	}
	_, err := uuid.Parse(s)
	if err != nil {
		fmt.Printf("testutil: uuid parse error: %v\n", err)
	}
	return err == nil
}

func (m *uuidMatcher) String() string {
	return "is UUID"
}

func IsUUID() gomock.Matcher {
	return &uuidMatcher{}
}
