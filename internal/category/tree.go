package category

import "shop/internal/models"

func buildTree(list []models.Category) []*models.CategoryNode {
	nodeMap := make(map[int]*models.CategoryNode)
	var tree []*models.CategoryNode

	for _, item := range list {
		nodeMap[item.ID] = &models.CategoryNode{
			Category: item,
			Children: []*models.CategoryNode{},
		}
	}

	for _, item := range list {
		node := nodeMap[item.ID]

		if item.ParentID != nil {
			parent, exists := nodeMap[*item.ParentID]
			if exists {
				parent.Children = append(parent.Children, node)
			}
		} else {
			tree = append(tree, node)
		}
	}

	return tree
}
