/**
* @file logger.go
* @author ligang
* @date 2016-02-04
 */

package golog

type simpleLogger struct {
	writer   IWriter
	formater IFormater

	glevel int
}

func NewSimpleLogger(writer IWriter, formater IFormater) *simpleLogger {
	return &simpleLogger{
		writer:   writer,
		formater: formater,

		glevel: LEVEL_INFO,
	}
}

func (s *simpleLogger) SetLogLevel(level int) *simpleLogger {
	_, ok := LogLevels[level]
	if ok {
		s.glevel = level
	}

	return s
}

func (s *simpleLogger) Debug(msg []byte) {
	s.Log(LEVEL_DEBUG, msg)
}

func (s *simpleLogger) Info(msg []byte) {
	s.Log(LEVEL_INFO, msg)
}

func (s *simpleLogger) Notice(msg []byte) {
	s.Log(LEVEL_NOTICE, msg)
}

func (s *simpleLogger) Warning(msg []byte) {
	s.Log(LEVEL_WARNING, msg)
}

func (s *simpleLogger) Error(msg []byte) {
	s.Log(LEVEL_ERROR, msg)
}

func (s *simpleLogger) Critical(msg []byte) {
	s.Log(LEVEL_CRITICAL, msg)
}

func (s *simpleLogger) Alert(msg []byte) {
	s.Log(LEVEL_ALERT, msg)
}

func (s *simpleLogger) Emergency(msg []byte) {
	s.Log(LEVEL_EMERGENCY, msg)
}

func (s *simpleLogger) Log(level int, msg []byte) error {
	if level > s.glevel {
		return nil
	}

	_, err := s.writer.Write(s.formater.Format(level, append(msg, '\n')))

	return err
}
