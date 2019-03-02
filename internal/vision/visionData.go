package vision

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
)

func micronizeBalls(balls []*sslproto.SSL_DetectionBall) (microBalls []*sslproto.SSL_Micro_DetectionBall) {
	microBalls = make([]*sslproto.SSL_Micro_DetectionBall, len(balls))
	for i, b := range balls {
		microBalls[i] = new(sslproto.SSL_Micro_DetectionBall)
		microBalls[i].X = b.X
		microBalls[i].Y = b.Y
	}
	return
}

func micronizeBots(robots []*sslproto.SSL_DetectionRobot) (microRobots []*sslproto.SSL_Micro_DetectionRobot) {
	microRobots = make([]*sslproto.SSL_Micro_DetectionRobot, len(robots))
	for i, r := range robots {
		microRobots[i] = new(sslproto.SSL_Micro_DetectionRobot)
		microRobots[i].RobotId = r.RobotId
		microRobots[i].X = r.X
		microRobots[i].Y = r.Y
		microRobots[i].Orientation = r.Orientation
	}
	return
}

func micronizeGeometry(geometry *sslproto.SSL_GeometryData) (microGeometry *sslproto.SSL_Micro_GeometryData) {
	microGeometry = new(sslproto.SSL_Micro_GeometryData)
	microGeometry.Field = new(sslproto.SSL_Micro_GeometryFieldSize)
	microGeometry.Field.BoundaryWidth = geometry.Field.BoundaryWidth
	microGeometry.Field.FieldLength = geometry.Field.FieldLength
	microGeometry.Field.FieldWidth = geometry.Field.FieldWidth
	microGeometry.Field.GoalDepth = geometry.Field.GoalDepth
	microGeometry.Field.GoalWidth = geometry.Field.GoalWidth
	microGeometry.Field.FieldLines = micronizeLines(geometry.Field.FieldLines)
	microGeometry.Field.FieldArcs = micronizeArcs(geometry.Field.FieldArcs)
	return
}
func micronizeLines(lines []*sslproto.SSL_FieldLineSegment) (microLines []*sslproto.SSL_Micro_FieldLineSegment) {
	microLines = make([]*sslproto.SSL_Micro_FieldLineSegment, len(lines))
	for i, r := range lines {
		microLines[i] = new(sslproto.SSL_Micro_FieldLineSegment)
		microLines[i].P1 = new(sslproto.Micro_Vector2F)
		microLines[i].P1.X = r.P1.X
		microLines[i].P1.Y = r.P1.Y
		microLines[i].P2 = new(sslproto.Micro_Vector2F)
		microLines[i].P2.X = r.P2.X
		microLines[i].P2.Y = r.P2.Y
	}
	return
}

func micronizeArcs(arcs []*sslproto.SSL_FieldCicularArc) (microArcs []*sslproto.SSL_Micro_FieldCicularArc) {
	microArcs = make([]*sslproto.SSL_Micro_FieldCicularArc, len(arcs))
	for i, r := range arcs {
		microArcs[i] = new(sslproto.SSL_Micro_FieldCicularArc)
		microArcs[i].Center = new(sslproto.Micro_Vector2F)
		microArcs[i].Center.X = r.Center.X
		microArcs[i].Center.Y = r.Center.Y
		microArcs[i].Radius = r.Radius
		microArcs[i].A1 = r.A1
		microArcs[i].A2 = r.A2
	}
	return
}
