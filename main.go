package main

import (
	"context"
	pb "createEvents/proto"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Schedule struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Class string `json:"class"`
}

var schedule []Schedule

type server struct {
	pb.UnimplementedParserServer
}

func (s *server) SendSchedule(ctx context.Context, req *pb.ScheduleRequest) (*pb.ScheduleResponse, error) {
	for _, entry := range req.Entries {
		fmt.Printf("Received schedule entry: Start=%s, End=%s, ClassName=%s\n", entry.Start, entry.End, entry.ClassName)
		schedule = append(schedule, Schedule{
			Start: entry.Start,
			End:   entry.End,
			Class: entry.ClassName,
		})
	}
	createEvent(schedule)
	return &pb.ScheduleResponse{Status: "OK"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterParserServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func createEvent(schedule []Schedule) {
	calendarID := ""
	keyFile := "credentials.json"

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile(keyFile))
	if err != nil {
		log.Fatalf("Failed to create Calendar service: %v", err)
	}
	for _, item := range schedule {
		if item.Class != "" {
			event := &calendar.Event{
				Summary: item.Class,
				Start: &calendar.EventDateTime{
					DateTime: item.Start,
					TimeZone: "Europe/Samara",
				},
				End: &calendar.EventDateTime{
					DateTime: item.End,
					TimeZone: "Europe/Samara",
				},
			}
			event, err = srv.Events.Insert(calendarID, event).Do()
			if err != nil {
				log.Fatalf("Error creating event: %v", err)
			}
			fmt.Printf("Event created: %s\n", event.HtmlLink)
		}
	}
}
