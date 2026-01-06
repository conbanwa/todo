package todo

import "sort"

// FilterAndSort filters the provided todos according to opts.Status and
// sorts them according to opts.SortBy and opts.SortOrder.
func FilterAndSort(in []Todo, opts ListOptions) []Todo {
	out := make([]Todo, 0, len(in))
	for _, v := range in {
		if opts.Status != "" && v.Status != opts.Status {
			continue
		}
		out = append(out, v)
	}

	cmp := func(i, j int) bool { return out[i].ID < out[j].ID }
	switch opts.SortBy {
	case "due_date":
		cmp = func(i, j int) bool { return out[i].DueDate.Before(out[j].DueDate) }
	case "status":
		cmp = func(i, j int) bool { return out[i].Status < out[j].Status }
	case "name":
		cmp = func(i, j int) bool { return out[i].Name < out[j].Name }
	}

	sort.SliceStable(out, func(i, j int) bool {
		if opts.SortOrder == "desc" {
			return !cmp(i, j)
		}
		return cmp(i, j)
	})

	return out
}
