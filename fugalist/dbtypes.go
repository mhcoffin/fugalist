package fugalist

import (
	"time"
)

type PaletteSoundId = string

type VstSoundId = string
type VstSound struct {
	Id       VstSoundId `firestore:"Id"`
	Name     string     `firestore:"name"`
	Midi     string     `firestore:"midi"`
	Stop     string     `firestore:"stop"`
	Dynamics string     `firestore:"dyn"`
}

type BranchId = string
type Branch struct {
	Id         BranchId
	Order      float64
	Condition  string
	VstSoundId string
	Length     float64
	Transpose  float64
}

type CompositeSoundId = string
type CompositeSound struct {
	Id       CompositeSoundId
	Name     string
	Branches map[BranchId]Branch
	Order    float64
}

type TechniqueId = string

type Technique struct {
	Id   TechniqueId `firestore:"id"`
	Name string      `firestore:"name"`
}

type AxisId = string

type Axis struct {
	Id         AxisId      `firestore:"id"`
	Name       string      `firestore:"name"`
	Techniques []Technique `firestore:"techniques"`
	SortOrder  float64     `firestore:"sortOrder"`
}

type ProjectId = string

type Metadata struct {
	ExpressionMapId string `firestore:"expressionMapId"`
	Plugins         string `firestore:"plugins"`
	Version         int    `firestore:"version"`
	Description     string `firestore:"description"`
	MiddleC         string
}

type Assignment struct {
	Sound string `firestore:"sound"`
}

type Tint = struct {
	Id    string `firestore:"id"`
	Order int    `firestore:"order"`
	Name  string `firestore:"name"`
	Midi  string `firestore:"midi"`
	Stop  string `firestore:"stop"`
}

type Project struct {
	ProjectId       string
	CreateTime      time.Time
	ModifyTime      time.Time
	Axes            map[string]Axis
	VstSounds       map[VstSoundId]*VstSound
	Tints           map[string]*Tint
	CompositeSounds map[CompositeSoundId]*CompositeSound
	Assignments     map[string]Assignment
	MiddleC         string
}

type AudioExample struct {
	Id    string
	Order float64
	Name  string
	URL   string
	Score string
}

type ProjectSummary struct {
	CreateTime        time.Time
	ModifyTime        time.Time
	ShareTime         time.Time `firestore:",serverTimestamp"`
	ProjectID         string
	Version           int
	Name              string
	Public            bool
	Description       string
	Plugins           string
	Vendor            string
	Instruments       string
	OtherTags         string
	Examples          map[string]AudioExample
	ExpressionMapURL  string
	ExpressionMapTime time.Time
}

type UserInfo struct {
	CanonicalDisplayName string
	CreationTime         time.Time `firestore:"serverTimestamp"`
	DisplayName          string
	Email                string
	PhotoURL             string
	Theme                string
}

type Share struct {
	ID              string
	CreateTime      time.Time `firestore:",serverTimestamp"`
	UID             string
	UserDisplayName string
	PhotoURL        string
	PID             string
	Summary         ProjectSummary
	Axes            []Axis
	Vendor          string
	Instruments     []string
	OtherTags       []string
	Tags            []string
	Superseded      bool
}
