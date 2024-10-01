package hogwarts

func GetAllCourses(prereqs map[string][]string) []string {
	set := make(map[string]struct{})
	for course, prereqList := range prereqs {
		set[course] = struct{}{}
		for _, prereq := range prereqList {
			set[prereq] = struct{}{}
		}
	}
	result := make([]string, 0, len(set))
	for course := range set {
		result = append(result, course)
	}
	return result
}

func GetCourseList(prereqs map[string][]string) []string {
	if len(prereqs) == 0 {
		return []string{}
	}

	courses := GetAllCourses(prereqs)
	result := make([]string, 0, len(courses))

	// 0 - not visited (WHITE)
	// 1 - in progress (GRAY)
	// 2 - completed (BLACK)
	state := make(map[string]int)

	var dfs func(course string)
	dfs = func(course string) {
		if state[course] == 2 { // ← упрощенная версия
			return
		}
		if state[course] == 1 {
			panic("cycle detected")
		}

		state[course] = 1
		for _, prereq := range prereqs[course] {
			dfs(prereq)
		}
		state[course] = 2
		result = append(result, course)
	}

	for _, course := range courses {
		dfs(course)
	}

	return result
}
