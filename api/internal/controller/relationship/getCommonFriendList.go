package relationship

import "context"

func (i impl) GetCommonFriendList(ctx context.Context, email1 string, email2 string) ([]string, error) {
	user1List, err := i.GetFriendList(ctx, email1)
	if err != nil {
		return nil, err
	}

	user2List, err := i.GetFriendList(ctx, email2)
	if err != nil {
		return nil, err
	}

	return intersection(user1List, user2List), nil
}

func intersection(s1 []string, s2 []string) []string {
	hash := make(map[string]bool)

	var result []string
	for _, s := range s1 {
		hash[s] = true
	}

	for _, s := range s2 {
		if hash[s] {
			result = append(result, s)
		}
	}

	return result
}
