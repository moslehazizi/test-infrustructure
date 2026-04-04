package entity

import (
	"control-panel-service/pkg"
	"time"

	"gorm.io/gorm"
)

type ScenarioStatus string

const (
	ScenarioStatusPending ScenarioStatus = "pending"
	ScenarioStatusRunning ScenarioStatus = "running"
	ScenarioStatusPaused  ScenarioStatus = "paused"
	ScenarioStatusStopped ScenarioStatus = "stopped"
	ScenarioStatusAborted ScenarioStatus = "aborted"
	ScenarioStatusSucceed ScenarioStatus = "succeed"
)

type TestScenario struct {
	ID                  uint64             `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt           time.Time          `gorm:"column:created_at"`
	UpdatedAt           time.Time          `gorm:"column:updated_at"`
	DeletedAt           *gorm.DeletedAt    `gorm:"column:deleted_at"`
	Name                string             `gorm:"column:name"`
	TestCategoryID      uint64             `gorm:"column:test_category_id"`
	TestCategory        *TestCategory      `gorm:"ForeignKey:TestCategoryID"`
	MotherServiceID     uint64             `gorm:"column:mother_service_id"`
	MotherService       *MotherService     `gorm:"ForeignKey:MotherServiceID"`
	Status              ScenarioStatus     `gorm:"column:status"`
	NumSteps            int64              `gorm:"column:num_steps"`
	MaxTestServiceCount *int64             `gorm:"column:max_test_service_count"`
	TestServiceConfig   *TestServiceConfig `gorm:"ForeignKey:TestScenarioID"`
	DeploymentNumber    int32              `gorm:"column:deployment_number"`
	StartedAt           *time.Time         `gorm:"column:started_at"`
	Editable            bool               `gorm:"column:editable;default:true;not null"`
	// سناریو را میتوان با افزایش تعداد ایجنت های تست به صورت خودکار مجدد اجرا کرد.
	// این فیلد تعداد ایجنت هایی که در هر اجرای سناریو از ابتدا لازم هست
	// که به تعداد مرحله قبلش اضافه شود را مشخص میکند.
	// برای مثال در هر اجرای کامل سناریو، ۵ ایجنت اضافه شود.
	// این به این معنا است که زمانیکه سناریو تا انتها رفت،
	// مجدد از ابتدا تکرار میشود و اینبار، تعداد ایجنت ۵ تا
	// بیشتر از اجرای قبلی خواهد بود.
	// نکته مهم اینکه فیلد اگزکیوشن آی دی در هر اجرای از ابتدا،
	// متفاوت خواهد بود و میتوان تمایز داد.
	IncreaseAgentNumber int64 `gorm:"column:increase_agent_number"`
	// این فیلد در کنار فیلد بالا معنا دار میشود.
	// این فیلد تعیین میکند که افزایش ایجنت ها چند بار رخ دهد.
	// برای مثال اگر این عدد ۱۰ باشد، یعنی ۱۰ بار کل سناریو رو تکرار کن
	// و هر بار به میزان تعریف شده در فیلد بالا، تعداد ایجنت ها را افزایش بده.
	//
	// فرض کنید:
	// IncreaseAgentNumber=5 & ExecNumMultiAgent=10 & MaxTestServiceCount = 5
	// در این حالت، ابتدا با ۵ ایجنت سناریو انجام میشود.
	// سپس با ۵ واحد افزایش، کل سناریو با ۱۰ ایجنت تکرار میشود.
	// این فرایند ۱۰ بار تکرار میشود و آخرین مرحله با ۵۵ ایجنت تست خواهد شد.
	ExecNumMultiAgent int64 `gorm:"column:execution_number_multi_agent"`
}

func (TestScenario) TableName() string {
	return "test_scenarios"
}

type TestScenarioPaginationRequest struct {
	Page    int
	PerPage int
}

func (ts *TestScenario) Validate(testCat *TestCategory) error {
	if ts.MaxTestServiceCount != nil && *ts.MaxTestServiceCount < 1 {
		return pkg.ErrMaxTestServiceCountLessThanOne
	}

	if testCat.HasNumSteps {
		if ts.NumSteps < 1 {
			return pkg.ErrNumStepsNotSet
		}
	} else {
		// if this field be false it means that we have a kind of test that has one step
		if ts.NumSteps != 1 {
			return pkg.ErrNumStepsShouldBeOne
		}
	}

	if testCat.HasMaxTestServiceCount {
		if ts.MaxTestServiceCount == nil {
			return pkg.ErrMaxTestServiceCountNotSet
		}
	} else {
		if ts.MaxTestServiceCount != nil {
			return pkg.ErrNoNeedMaxTestServiceCount
		}
	}

	if ts.IncreaseAgentNumber < 0 {
		return pkg.ErrIncreaseAgentNumNotBeNegative
	}
	if ts.ExecNumMultiAgent < 1 {
		return pkg.ErrExecNumMultiAgentShouldBePositive
	}

	return nil
}
