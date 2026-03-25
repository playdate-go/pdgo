package utils

import "fmt"

const PlayDateSDKPathEnvVar = "PLAYDATE_SDK_PATH"

var ErrPlayDateSDKPathNotSet = fmt.Errorf("%s isn't set and no fallback value is available, please set env var to SDK path", PlayDateSDKPathEnvVar)
var ErrNoFallbackValue = fmt.Errorf("no fallback value")
