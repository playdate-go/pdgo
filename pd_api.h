// IMPORTANT: Set CGO_CFLAGS to point to your Playdate SDK:
//   export CGO_CFLAGS="-I$HOME/Developer/PlaydateSDK/C_API -DTARGET_EXTENSION=1"
//
// Or use pdgoc which handles SDK paths automatically.

#ifndef PDGO_API_H
#define PDGO_API_H

#include <stdint.h>
#include <stdlib.h>

#include "pd_api/pd_api_file.h"
#include "pd_api/pd_api_gfx.h"
#include "pd_api/pd_api_sys.h"
#include "pd_api/pd_api_lua.h"
#include "pd_api/pd_api_json.h"
#include "pd_api/pd_api_sprite.h"
#include "pd_api/pd_api_sound.h"
#include "pd_api/pd_api_display.h"
#include "pd_api/pd_api_scoreboards.h"
#include "pd_api/pd_api_network.h"

typedef struct PlaydateAPI PlaydateAPI;

struct PlaydateAPI
{
    const struct playdate_sys* system;
    const struct playdate_file* file;
    const struct playdate_graphics* graphics;
    const struct playdate_sprite* sprite;
    const struct playdate_display* display;
    const struct playdate_sound* sound;
    const struct playdate_lua* lua;
    const struct playdate_json* json;
    const struct playdate_scoreboards* scoreboards;
    const struct playdate_network* network;
};

#if TARGET_EXTENSION
typedef enum
{
    kEventInit,
    kEventInitLua,
    kEventLock,
    kEventUnlock,
    kEventPause,
    kEventResume,
    kEventTerminate,
    kEventKeyPressed,
    kEventKeyReleased,
    kEventLowPower,
    kEventMirrorStarted,
    kEventMirrorEnded
} PDSystemEvent;
#endif

#ifdef _WINDLL
__declspec(dllexport)
#endif
int eventHandler(PlaydateAPI* playdate, PDSystemEvent event, uint32_t arg);

#endif 
