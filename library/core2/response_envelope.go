package core2

import "gitlab.int.magneato.site/dungar/prototype/internal/utils"

// ResponseEnvelope is a list of response to a given target
type ResponseEnvelope struct {
	Message   *IncomingMessage
	Responses []*Response
}

// DeDuplicate responses
func (re *ResponseEnvelope) DeDuplicate(driver ProtocolDriver) {
	newResponses := make([]*Response, 0)
	tmp := make([]string, 0)

	for _, rsp := range re.Responses {
		rspTxt := rsp.Format(driver, re.Message)

		if utils.StringInSlice(rspTxt, tmp) {
			continue
		}

		newResponses = append(newResponses, rsp)
		tmp = append(tmp, rspTxt)
	}

	re.Responses = newResponses
}
