package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/mirzaakhena/gogen/gogen"

	"github.com/mirzaakhena/seqode/model"
)

func main() {
	content, err := ioutil.ReadFile("sample/sample7.wsd")
	if err != nil {
		log.Fatal(err)
	}

	rightArrow := regexp.MustCompile(`[\w\d]+(\s)*(->)(\s)*[\w\d]+(\s)*:(\s)*`)
	rightReturnArrow := regexp.MustCompile(`[\w\d]+(\s)*(-->)(\s)*[\w\d]+(\s)*:(\s)*`)

	root := model.Root{
		Participants: map[string]*model.Participant{},
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		row := strings.TrimSpace(scanner.Text())

		if len(row) == 0 {
			continue
		} else //

		if strings.HasPrefix(row, "@startuml") {
			root.Name = row[10:]
		} else //

		if strings.HasPrefix(row, "@enduml") {
			break
		} else //

		if strings.HasPrefix(row, "participant") {
			participantName := row[12:]
			root.Participants[participantName] = &model.Participant{
				Name: participantName,
			}
		} else //

		if rightArrow.Match([]byte(row)) {

			leftParticipantName, rightParticipantName, usecaseName := GetInteraction(row, "->")

			if len(usecaseName) == 0 {
				log.Fatal("usecase name must not empty")
			}

			rightParticipant, ok := root.Participants[rightParticipantName]
			if !ok {
				continue
			}

			if rightParticipant.Usecases == nil {
				rightParticipant.Usecases = map[string]*model.Usecase{}
			}

			rightParticipant.Usecases[usecaseName] = &model.Usecase{
				Name:        usecaseName,
				Participant: rightParticipant,
			}

			rightParticipant.CurrentUsecase = rightParticipant.Usecases[usecaseName]

			leftParticipant := root.Participants[leftParticipantName]

			if leftParticipant.CurrentUsecase != nil {
				leftParticipant.CurrentUsecase.Outports = append(leftParticipant.CurrentUsecase.Outports, rightParticipant.Usecases[usecaseName])

				if leftParticipant.CurrentParticipantOutport == nil {
					leftParticipant.CurrentParticipantOutport = rightParticipant
				} else if leftParticipant.CurrentParticipantOutport.CurrentUsecase != nil {
					log.Fatalf("Left Participant's Outport %s.%s is not closed yet", leftParticipantName, leftParticipant.CurrentParticipantOutport.CurrentUsecase.Name)
				}
			}

		} else //

		if rightReturnArrow.Match([]byte(row)) {

			leftParticipantName, _, _ := GetInteraction(row, "-->")

			leftParticipant, ok := root.Participants[leftParticipantName]
			if !ok {
				continue
			}

			if leftParticipant.CurrentParticipantOutport != nil {
				if leftParticipant.CurrentParticipantOutport.CurrentUsecase != nil {
					log.Fatalf("Right Participant %v.%v is not closed yet", leftParticipant.Name, leftParticipant.CurrentParticipantOutport.CurrentUsecase.Name)
				}

			}

			leftParticipant.CurrentUsecase = nil
			leftParticipant.CurrentParticipantOutport = nil

		}

	}

	fmt.Println()

	for _, x := range root.Participants {
		printAll(x)
		gogen.GenerateInit(gogen.InitRequest{FolderPath: fmt.Sprintf("hello/%s", x.Name)})
		for _, u := range x.Usecases {

			methods := []string{}
			for _, o := range u.Outports {
				methods = append(methods, o.Name)
			}

			gogen.GenerateUsecase(gogen.UsecaseRequest{
				UsecaseType:    "command",
				UsecaseName:    u.Name,
				OutportMethods: methods,
				FolderPath:     fmt.Sprintf("hello/%s", x.Name),
			})

		}
	}

}

func printAll(x *model.Participant) {

	fmt.Printf("participant: %v\n", x.Name)
	for _, u := range x.Usecases {
		fmt.Printf(" have usecase %v\n", u.Name)
		for _, o := range u.Outports {
			fmt.Printf("  have outport %v.%v\n", o.Participant.Name, o.Name)
		}
	}
	fmt.Println()

}

func GetInteraction(row, arrow string) (string, string, string) {

	indexArrow := strings.Index(row, arrow)
	colonIndex := strings.Index(row, ":")

	leftParticipantName := strings.TrimSpace(row[:indexArrow])
	rightParticipantName := strings.TrimSpace(row[indexArrow+len(arrow) : colonIndex])
	usecaseName := strings.TrimSpace(row[colonIndex+1:])

	return leftParticipantName, rightParticipantName, usecaseName
}
