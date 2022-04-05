package response

import "github.com/kubeedge/mappers-go/mapper-sdk-go/internal/common"
import "net/http"

func CodeMapping(kind common.ErrKind) int {
	if kind == "" {
		return 200
	}
	switch kind {
	case common.KindServerError:
		return http.StatusInternalServerError
	case common.KindEntityDoesNotExist:
		return http.StatusRequestedRangeNotSatisfiable
	case common.KindInvalidId:
		return http.StatusBadGateway
	case common.KindServiceUnavailable:
		return http.StatusServiceUnavailable
	case common.KindServiceLocked:
		return http.StatusLocked
	case common.KindNotImplemented:
		return http.StatusNotImplemented
	case common.KindRangeNotSatisfiable:
		return http.StatusRequestedRangeNotSatisfiable
	case common.KindOverflowError, common.KindNaNError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
