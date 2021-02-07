package types

type RLog []*LogTuple

func IntersectionInt(s1, s2 []int) (inter []int) {
	hash := make(map[int]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = RemoveDupsInt(inter)
	return

}
func RemoveDupsInt(elements []int) (nodups []int) {
	encountered := make(map[int]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func AppendIfMissingInt(slice []int, i int) []int {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
func AppendIfMissingRequest(slice []*Request, i *Request) []*Request {
	for _, ele := range slice {
		if ele.Equals(i) {
			return slice
		}
	}
	return append(slice, i)
}

func AppendIfMissingAcceptedRequest(slice []*AcceptedRequest, i *AcceptedRequest) []*AcceptedRequest {
	for _, ele := range slice {
		if ele.Equals(i) {
			return slice
		}
	}
	return append(slice, i)
}

func AppendIfMissingRequestStatus(slice []*RequestStatus, i *AcceptedRequest) []*RequestStatus {
	for _, ele := range slice {
		if ele.Req.Equals(i) {
			return slice
		}
	}
	return append(slice, &RequestStatus{Req: i})
}

func (tuples1 RLog) CommonPrefix(tuples2 RLog) RLog {
	var logs RLog = make([]*LogTuple, 0)
	var min RLog
	if len(tuples1) < len(tuples2) {
		min = tuples1
	} else {
		min = tuples2
	}
	for i := range min {
		if tuples1[i].Equals(tuples2[i]) {
			logs = append(logs, tuples1[i])
		}
	}
	return logs
}

func ExcludeRequests(src []*Request, target []*Request) []*Request {
	out := make([]*Request, 0)
	for _, s := range src {
		flag := true
		for _, r := range target {
			if s.Equals(r) {
				flag = false
				break
			}
		}
		if flag {
			out = append(out, s)
		}
	}
	return out
}

func FilterRequests(src []*Request, filterFn func(request *Request) bool) []*Request {
	out := make([]*Request, 0)
	for _, s := range src {
		if !filterFn(s) {
			out = append(out, s)
		}
	}
	return out
}

func FilterAcceptedRequests(src []*AcceptedRequest, filterFn func(request *AcceptedRequest) bool) []*AcceptedRequest {
	out := make([]*AcceptedRequest, 0)
	for _, s := range src {
		if !filterFn(s) {
			out = append(out, s)
		}
	}
	return out
}

func FilterRequestStatus(src []*RequestStatus,
	filterFn func(rs *RequestStatus) bool,
) []*RequestStatus {
	out := make([]*RequestStatus, 0)
	for _, s := range src {
		if !filterFn(s) {
			out = append(out, s)
		}
	}
	return out
}

func RLogEquals(rl1 []*LogTuple, rl2 []*LogTuple) bool {
	if len(rl1) != len(rl2) {
		return false
	}
	for i, log := range rl1 {
		if !log.Equals(rl2[i]) {
			return false
		}
	}
	return true
}

func LastReqEquals(lr1 []*RequestReply, lr2 []*RequestReply) bool {
	if len(lr1) != len(lr2) {
		return false
	}
	for i, log := range lr1 {
		if !log.Equals(lr2[i]) {
			return false
		}
	}
	return true
}
