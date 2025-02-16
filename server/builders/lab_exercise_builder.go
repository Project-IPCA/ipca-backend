package builders

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type LabExerciseBuilder struct {
	ExerciseID             uuid.UUID
	ChapterID              *uuid.UUID
	Level                  *string
	Name                   *string
	Content                *string
	Testcase               string
	Sourcecode             *string
	FullMark               int
	AddedDate              time.Time
	LastUpdate             *time.Time
	UserDefinedConstraints *json.RawMessage
	SuggestedConstraints   *json.RawMessage
	AddedBy                *string
	CreatedBy              *uuid.UUID
	Language               *string
}

func NewLabExerciseBuilder() *LabExerciseBuilder {
	return &LabExerciseBuilder{
		ExerciseID: uuid.New(),
		Testcase:   "NO_INPUT",
		FullMark:   10,
		AddedDate:  time.Now(),
	}
}

func (b *LabExerciseBuilder) SetExerciseID(id uuid.UUID) *LabExerciseBuilder {
	b.ExerciseID = id
	return b
}

func (b *LabExerciseBuilder) SetChapterID(id *uuid.UUID) *LabExerciseBuilder {
	b.ChapterID = id
	return b
}

func (b *LabExerciseBuilder) SetLevel(level *string) *LabExerciseBuilder {
	b.Level = level
	return b
}

func (b *LabExerciseBuilder) SetName(name *string) *LabExerciseBuilder {
	b.Name = name
	return b
}

func (b *LabExerciseBuilder) SetContent(content *string) *LabExerciseBuilder {
	b.Content = content
	return b
}

func (b *LabExerciseBuilder) SetTestcase(testcase string) *LabExerciseBuilder {
	b.Testcase = testcase
	return b
}

func (b *LabExerciseBuilder) SetSourcecode(sourcecode *string) *LabExerciseBuilder {
	b.Sourcecode = sourcecode
	return b
}

func (b *LabExerciseBuilder) SetFullMark(fullMark int) *LabExerciseBuilder {
	b.FullMark = fullMark
	return b
}

func (b *LabExerciseBuilder) SetAddedDate(addedDate time.Time) *LabExerciseBuilder {
	b.AddedDate = addedDate
	return b
}

func (b *LabExerciseBuilder) SetLastUpdate(lastUpdate *time.Time) *LabExerciseBuilder {
	b.LastUpdate = lastUpdate
	return b
}

func (b *LabExerciseBuilder) SetUserDefinedConstraints(
	constraints *json.RawMessage,
) *LabExerciseBuilder {
	b.UserDefinedConstraints = constraints
	return b
}

func (b *LabExerciseBuilder) SetSuggestedConstraints(
	constraints *json.RawMessage,
) *LabExerciseBuilder {
	b.SuggestedConstraints = constraints
	return b
}

func (b *LabExerciseBuilder) SetAddedBy(addedBy *string) *LabExerciseBuilder {
	b.AddedBy = addedBy
	return b
}

func (b *LabExerciseBuilder) SetCreatedBy(createdBy *uuid.UUID) *LabExerciseBuilder {
	b.CreatedBy = createdBy
	return b
}

func (b *LabExerciseBuilder) SetLanguage(language *string) *LabExerciseBuilder {
	b.Language = language
	return b
}

func (b *LabExerciseBuilder) Build() models.LabExercise {
	return models.LabExercise{
		ExerciseID:             b.ExerciseID,
		ChapterID:              b.ChapterID,
		Level:                  b.Level,
		Name:                   b.Name,
		Content:                b.Content,
		Testcase:               b.Testcase,
		Sourcecode:             b.Sourcecode,
		FullMark:               b.FullMark,
		AddedDate:              b.AddedDate,
		LastUpdate:             b.LastUpdate,
		UserDefinedConstraints: b.UserDefinedConstraints,
		SuggestedConstraints:   b.SuggestedConstraints,
		AddedBy:                b.AddedBy,
		CreatedBy:              b.CreatedBy,
		Language:               b.Language,
	}
}
