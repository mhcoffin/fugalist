package fugalist

import "time"

type PaletteSoundId = string

type SoundId = string

type PigmentId = string
type Pigment struct {
	PigmentId string `firestore:"Id"`
	Name      string `firestore:"name"`
	Midi      string `firestore:"midi"`
	Stop      string `firestore:"stop"`
	Dynamics  string `firestore:"dyn"`
}

type BranchId = string
type Branch struct {
	Id        BranchId
	Condition string
	Length    float64
	Order     float32
	Transpose float64
	Pigment   string
}

type ColorId = string
type Color struct {
	ColorId  string
	Branches map[BranchId]Branch
}

type TechniqueId = string

type Technique struct {
	Id   TechniqueId `firestore:"id"`
	Name string      `firestore:"name"`
	Midi string      `firestore:"midi"`
}

type AxisId = string

type Axis struct {
	Id         AxisId      `firestore:"id"`
	Name       string      `firestore:"name"`
	AddOn      bool        `firestore:"addOn"`
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
	ProjectId   string
	Owner       string
	Version     int
	Public      bool
	CreateTime  time.Time
	ModifyTime  time.Time
	Axes        map[string]Axis
	Pigments    map[PigmentId]*Pigment
	Palette     map[ColorId]*Color
	Tints       map[string]*Tint
	Assignments map[string]Assignment
	URL         string
	MiddleC		string
}

// export interface Project {
// ProjectID: ProjectId;
// Version: number;
// CreateTime: Timestamp;
// ModifyTime: Timestamp;
// Name: string;
// MiddleC: string
// Axes: Record<string, Axis>;
// Pigments: Record<string, Pigment>;
// Tints: Record<string, Tint>
// Palette: Record<string, Color>
// Assignments: Assignments;
// URL: string | null;
// }

type ProjectSummary struct {
	CreateTime       time.Time
	ProjectID        string
	Version          int
	Name             string
	Public           bool
	Description      string
	Plugins          string
	ExpressionMapURL string
}

type UserInfo struct {
	Projects    map[string]ProjectSummary `firestore:"Projects"`
	Preferences map[string]string         `firestore:"Preferences"`
}