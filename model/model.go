package model

type Root struct {
	Name         string
	Participants map[string]*Participant
}

type Participant struct {
	Name                      string
	CurrentUsecase            *Usecase
	CurrentParticipantOutport *Participant
	Usecases                  map[string]*Usecase
}

type Usecase struct {
	Name        string
	Participant *Participant
	Outports    []*Usecase
}
