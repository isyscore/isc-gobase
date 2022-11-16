package matcher

import (
	"github.com/isyscore/isc-gobase/constants"
	"reflect"
	"strconv"
	"strings"
)

func CollectAccept(objectTypeFullName string, _ reflect.Kind, objectFieldName string, tagName string, subCondition string, errMsg string) {
	if constants.Accept != tagName {
		return
	}

	accept, err := strconv.ParseBool(strings.TrimSpace(subCondition))
	if err != nil {
		return
	}
	addMatcher(objectTypeFullName, objectFieldName, nil, errMsg, accept)
}
