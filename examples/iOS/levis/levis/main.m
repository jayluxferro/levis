//
//  main.m
//  levis
//
//  Created by Jay on 10/21/21.
//

#import <UIKit/UIKit.h>
#import "Ioslevis/Ioslevis.h"
#define LOCAL_SERVER "0.0.0.0:5681"

void initServer(int argc, char *arg[]){
    NSString *endpoint = [NSString stringWithUTF8String:LOCAL_SERVER];
    if(strcmp(arg[1], "t") == 0){
        IoslevisStartTCPServer(endpoint);
    }else {
        IoslevisStartUDPServer(endpoint);
    }
}

void initClient(int argc, char *arg[]){
    if (argc >= 4){
        NSString *endpoint = [NSString stringWithCString:arg[3] encoding:NSUTF8StringEncoding];
        if(argc == 4){
            // no payload passed
            NSString *payload = @"";
            if(strcmp(arg[1], "t") == 0){
                IoslevisStartTCPClient(endpoint, payload);
            }else{
                IoslevisStartUDPClient(endpoint, payload);
            }
        }else {
            NSString *payload = [NSString stringWithCString:arg[4] encoding:NSUTF8StringEncoding];
            if(strcmp(arg[1], "t") == 0){
                IoslevisStartTCPClient(endpoint, payload);
            }else{
                IoslevisStartUDPClient(endpoint, payload);
            }
        }
    }
}


int main(int argc, char *arg[]) {
    @autoreleasepool {
        // get cmd line args
        // ./levis <proto|| tcp/udp> <c/s> ||[c => endpoint, => payload]
        if (argc < 3){
            return EXIT_FAILURE;
        }
        
        if (argc == 3){
            initServer(argc, arg);
        }else{
            initClient(argc, arg);
        }
    }
    return 0;
}

