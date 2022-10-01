package sql

const (
	MySqlPrefixPattern  = `^\d{13}[a-zA-Z0-9\-._]+$`
	MySqlPostfixPattern = `^[a-zA-Z0-9\-._]+\d{13}$`

	MSSQLPrefixPattern  = `[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]%`
	MSSQLPostfixPattern = `%[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]`
)
