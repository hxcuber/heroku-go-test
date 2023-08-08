package relationship

import "context"

func (i impl) GetCommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error) {
	user1Friends, err := i.GetFriends(ctx, email1)
	if err != nil {
		return nil, err
	}

	user2Friends, err := i.GetFriends(ctx, email2)
	if err != nil {
		return nil, err
	}

	return intersection(user1Friends, user2Friends), nil
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
