package gonmap

type ProbeList []string

var emptyProbeList []string

func (pl ProbeList) exist(ProbeName string) bool {
	for _, name := range pl {
		if name == ProbeName {
			return true
		}
	}
	return false
}

func (p ProbeList) removeDuplicate() ProbeList {
	res := make([]string, 0, len(p))
	temp := map[string]struct{}{}
	for _, item := range p {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			res = append(res, item)
		}
	}
	return res
}
