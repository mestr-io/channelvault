package app_errors
import 
(
	"fmt"
	"log"
	"runtime"
	"errors"
)

// Fatal logs a fatal error message and exits the application.
func Fatal(format string, v ...any) {
	where := getLineAndNo()
	s := fmt.Sprintf(format, v...)
	log.Fatalf("%s FATAL: %s", where, s)
}

// Error formats an error message with the file and line number where it was called.
func Error(format string, v ...any) (err error) {
	where := getLineAndNo()
	s := fmt.Sprintf(format, v...)
	err = fmt.Errorf("%s %s", where, s)

	return 
}

// Finally executes a function and captures any error, returning it along with the previous error.
func Finally(f func() error, previous error) (err error){	
	if previous != nil {
		err = previous
	}

	e := f()

	if e != nil {
		if previous != nil {
			err = errors.Join(previous, Error("%s",e.Error()))
		} else {
			err = e
		}		
	}

	return 	
}

func getLineAndNo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	
	return fmt.Sprintf("%s:%d", file, line)
}
