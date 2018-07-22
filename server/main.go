

//go:generate protoc -I ./helloworld --go_out=plugins=grpc:./helloworld ./helloworld/helloworld.proto

package main

import (
	"log"
	"net"
	"os/exec"
	"fmt"
	"unsafe"
	pb "github.com/kjunichi/gopheron-macos/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#include <stdio.h>
#include <stdlib.h>

char *getFinderIconStr() {
    //NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];
    NSString *path = @"/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/FinderIcon.icns";

    NSArray *imageReps = [NSBitmapImageRep imageRepsWithContentsOfFile:path];
    //NSInteger width = 0;
    //NSInteger height = 0;

    NSImageRep *imageRep = [imageReps objectAtIndex: 1];
    //NSLog(@"width = %ld",[imageRep pixelsWide]);
    //NSLog(@"height = %ld",[imageRep pixelsHigh]);
    
    NSData *data = [imageRep representationUsingType:NSPNGFileType properties : nil];
    NSString *pngstr = [data base64EncodedStringWithOptions:0 ];
    int len = [pngstr length];
    char *buf = (char*)malloc(sizeof(char)*(len+1));
    //NSLog(@"urldata = %@",pngstr);
    buf = [pngstr UTF8String];
    //printf("%s\n",buf);
    //[pool drain];
    return buf;
}
*/
import "C"

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	out, err := exec.Command("osx-cpu-temp").Output()
	if err != nil {
		log.Fatalf("Exec fail: %v", err)
	}
	// read finder icns
	tmpStr := C.getFinderIconStr()
	pngDataUrl := C.GoString(tmpStr)
	msg := fmt.Sprintf("<img width=\"128\" src=\"data:image/png;base64,%s\"><span style=\"font-size:larger;font-family: 'Segoe UI Emoji';\"><div style=\"font-size: 48px;wrap\">%s</div></span>",pngDataUrl,string(out))
	C.free(unsafe.Pointer(tmpStr))

	return &pb.HelloReply{Message: msg}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
