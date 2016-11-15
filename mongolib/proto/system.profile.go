package proto

import "time"

type ExecStats struct {
	ExecutionTimeMillisEstimate float64 `bson:"executionTimeMillisEstimate"`
	IsEOF                       float64 `bson:"isEOF"`
	NReturned                   float64 `bson:"nReturned"`
	NeedTime                    float64 `bson:"needTime"`
	RestoreState                float64 `bson:"restoreState"`
	Works                       float64 `bson:"works"`
	DocsExamined                float64 `bson:"docsExamined"`
	Direction                   string  `bson:"direction"`
	Invalidates                 float64 `bson:"invalidates"`
	NeedYield                   float64 `bson:"needYield"`
	SaveState                   float64 `bson:"saveState"`
	Stage                       string  `bson:"stage"`
	Advanced                    float64 `bson:"advanced"`
}

type MMAPV1Journal struct {
	AcquireCount AcquireCount `bson:"acquireCount"`
}

type Collection struct {
	AcquireCount AcquireCount `bson:"acquireCount"`
}

type SystemProfile struct {
	Query          map[string]interface{} `bson:"query"`
	Ts             time.Time              `bson:"ts"`
	Client         string                 `bson:"client"`
	Cursorid       float64                `bson:"cursorid"`
	ExecStats      ExecStats              `bson:"execStats"`
	Ns             string                 `bson:"ns"`
	Op             string                 `bson:"op"`
	WriteConflicts float64                `bson:"writeConflicts"`
	KeyUpdates     float64                `bson:"keyUpdates"`
	KeysExamined   float64                `bson:"keysExamined"`
	Locks          Locks                  `bson:"locks"`
	Nreturned      int                    `bson:"nreturned"`
	ResponseLength float64                `bson:"responseLength"`
	DocsExamined   int                    `bson:"docsExamined"`
	Millis         float64                `bson:"millis"`
	NumYield       float64                `bson:"numYield"`
	User           string                 `bson:"user"`
}
