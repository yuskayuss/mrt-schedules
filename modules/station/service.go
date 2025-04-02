package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/yuskayuss/mrt-schedules/common/client"
)

type Service interface {
	GetAllStation() ([]StationResponse, error)
	CheckSchedulesByStation(id string) ([]ScheduleResponse, error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStation() ([]StationResponse, error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return nil, err
	}

	var stations []Station
	if err := json.Unmarshal(byteResponse, &stations); err != nil {
		return nil, err
	}

	var response []StationResponse
	for _, item := range stations {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return response, nil
}

func (s *service) CheckSchedulesByStation(id string) ([]ScheduleResponse, error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return nil, err
	}

	var schedules []Schedule
	if err := json.Unmarshal(byteResponse, &schedules); err != nil {
		return nil, err
	}

	var scheduleSelected Schedule
	for _, item := range schedules {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		return nil, errors.New("Station not found")
	}

	return ConvertDataToResponse(scheduleSelected)
}

func ConvertDataToResponse(schedule Schedule) ([]ScheduleResponse, error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleLebakBulus)
	if err != nil {
		return nil, err
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleBundaranHI)
	if err != nil {
		return nil, err
	}

	var response []ScheduleResponse
	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	return response, nil
}

func ConvertScheduleToTimeFormat(schedule string) ([]time.Time, error) {
	schedules := strings.Split(schedule, ",")
	var response []time.Time

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}

		parsedTime, err := time.Parse("15:04", trimmedTime)
		if err != nil {
			return nil, errors.New("invalid time format: " + trimmedTime)
		}

		response = append(response, parsedTime)
	}

	return response, nil
}
