package vegvesen

import (
	"encoding/json"
	"strings"
)

type VehicleData struct {
	Brand                   string
	Color                   string
	Model                   string
	PersonalizedPlateNumber string
	RegistrationNumber      string
	VIN                     string
}

func (v *VehicleData) UnmarshalJSON(bytes []byte) error {
	var vehicleDataResponse struct {
		PersonalizedPlateNumber string `json:"personligKjennemerke"`
		RegistrationNumber      string `json:"kjennemerke"`
		Technical               struct {
			Model   string `json:"handelsbetegnelse"`
			Brand   string `json:"merke"`
			Chassis struct {
				Color string `json:"farge"`
			} `json:"karosseri"`
		} `json:"tekniskKjoretoy"`
		VIN string `json:"understellsnummer"`
	}

	if err := json.Unmarshal(bytes, &vehicleDataResponse); err != nil {
		return err
	}

	v.RegistrationNumber = vehicleDataResponse.RegistrationNumber
	v.PersonalizedPlateNumber = vehicleDataResponse.PersonalizedPlateNumber
	v.VIN = vehicleDataResponse.VIN
	v.Model = vehicleDataResponse.Technical.Model
	v.Brand = titleCase(vehicleDataResponse.Technical.Brand)
	v.Color = vehicleDataResponse.Technical.Chassis.Color
	return nil
}

func titleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

