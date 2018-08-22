package main

import "github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"

type Team struct {
	Name            string
	Goals           int
	RedCards        int
	YellowCards     int
	YellowCardTimes []uint32
	Timeouts        int
	TimeoutTime     int
}

type Stage struct {
	Name     string
	TimeLeft int
}

type Command struct {
	Name string
}

type Originator struct {
	Team  string
	BotId int
}

type GameEvent struct {
	Type       string
	Originator Originator
	Message    string
}

type Referee struct {
	Stage      Stage
	Command    Command
	TeamYellow Team
	TeamBlue   Team
	GameEvent  GameEvent
}

var referee = Referee{Stage{"NORMAL_FIRST_HALF_PRE", 1000}, Command{"HALT"},
	Team{"yellow", 5, 0, 1, []uint32{5000}, 1, 10000},
	Team{"blue", 1, 1, 3, []uint32{}, 2, 20000},
	GameEvent{"UNKNOWN", Originator{"UNKNOWN", -1}, "Custom message"}}

func saveRefereeMessageFields(message *sslproto.SSL_Referee) {
	referee.Stage.Name = message.Stage.String()
	if message.StageTimeLeft != nil {
		referee.Stage.TimeLeft = int(*message.StageTimeLeft)
	} else {
		referee.Stage.TimeLeft = 0
	}
	referee.Command.Name = message.Command.String()
	referee.TeamYellow = mapTeam(message.Yellow)
	referee.TeamBlue = mapTeam(message.Blue)
	if message.GameEvent != nil {
		if message.GameEvent.Originator != nil {
			referee.GameEvent.Originator.Team = message.GameEvent.Originator.Team.String()
			if message.GameEvent.Originator.BotId == nil {
				referee.GameEvent.Originator.BotId = -1
			} else {
				referee.GameEvent.Originator.BotId = int(*message.GameEvent.Originator.BotId)
			}
		} else {
			referee.GameEvent.Originator.Team = "TEAM_UNKNOWN"
			referee.GameEvent.Originator.BotId = -1
		}
		if message.GameEvent.Message == nil {
			referee.GameEvent.Message = ""
		} else {
			referee.GameEvent.Message = *message.GameEvent.Message
		}
		referee.GameEvent.Type = message.GameEvent.GameEventType.String()
	}
}

func mapTeam(teamInfo *sslproto.SSL_Referee_TeamInfo) (team Team) {
	team.Name = *teamInfo.Name
	team.Goals = int(*teamInfo.Score)
	team.YellowCards = int(*teamInfo.YellowCards)
	team.RedCards = int(*teamInfo.RedCards)
	team.YellowCardTimes = teamInfo.YellowCardTimes
	team.Timeouts = int(*teamInfo.Timeouts)
	team.TimeoutTime = int(*teamInfo.TimeoutTime)
	return
}
