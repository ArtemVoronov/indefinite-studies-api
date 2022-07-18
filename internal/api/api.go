package api

const (
	DUPLICATE_FOUND                      string = "DUPLICATE_FOUND"
	DONE                                 string = "DONE"
	DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN  string = "DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN"
	DELETE_VIA_POST_REQUEST_IS_FODBIDDEN string = "DELETE_VIA_POST_REQUEST_IS_FODBIDDEN"
	PAGE_NOT_FOUND                       string = "404 page not found"

	ERROR_MESSAGE_PARSING_BODY_JSON string = "Error during parsing of HTTP request body. Please check it format correctness: missed brackets, double quotes, commas, matching of names and data types and etc"
	ERROR_ID_WRONG_FORMAT           string = "Wrong ID format. Expected number"
)
