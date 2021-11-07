package error

import "fmt"

type ResourseError struct{
	Url string
	Err error
}

func (err *ResourseError) Error() string {
	return fmt.Sprintf("Error %s on url %s", err.Err, err.Url)
}