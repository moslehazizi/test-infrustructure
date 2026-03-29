package entity

import "time"

const (
	LOAD        string = "load"
	SMOKE       string = "smoke"
	SOAK        string = "soak"
	PEAK        string = "peak"
	SPIKE       string = "spike"
	SCALABILITY string = "scalability"
	STRESS      string = "stress"
	RECOVERY    string = "recovery"
)

type TestCategory struct {
	ID                     uint64 `gorm:"primarykey"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Name                   string
	Label                  string
	HasMaxTestServiceCount bool
	HasNumSteps            bool
	Active                 bool
}
