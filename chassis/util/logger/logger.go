package logger




type Logger interface {
	Error(format string) 
	Debug(format string) 
	Infor(format string) 
	Warning(format string) 
	Fatal(format string) 
}


type LogSearching interface {
	Find(index string, format string, start *string, end *string) []string;
}