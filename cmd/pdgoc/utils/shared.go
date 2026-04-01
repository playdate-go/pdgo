package utils

import "fmt"

const PlaydateSDKPathEnvVar = "PLAYDATE_SDK_PATH"

var ErrPlaydateSDKPathNotSet = fmt.Errorf("%s isn't set and no fallback value is available, please set env var to SDK path", PlaydateSDKPathEnvVar)
var ErrNoFallbackValue = fmt.Errorf("no fallback value")
