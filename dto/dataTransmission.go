package dto

type DataTransmission struct {
	Name	 string
	Data	 interface{}
}

func DTO(name string, data interface{}) DataTransmission {
	return DataTransmission{Name: name, Data: data}
}
