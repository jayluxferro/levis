// Objective-C API for talking to ioslevis Go package.
//   gobind -lang=objc ioslevis
//
// File is generated by gobind. Do not edit.

#ifndef __Ioslevis_H__
#define __Ioslevis_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


FOUNDATION_EXPORT void IoslevisStartTCPClient(NSString* _Nullable endpoint, NSString* _Nullable payload);

/**
 * *** TCP ***
 */
FOUNDATION_EXPORT void IoslevisStartTCPServer(NSString* _Nullable endpoint);

FOUNDATION_EXPORT void IoslevisStartUDPClient(NSString* _Nullable endpoint, NSString* _Nullable payload);

/**
 * ** UDP **
 */
FOUNDATION_EXPORT void IoslevisStartUDPServer(NSString* _Nullable endpoint);

#endif
