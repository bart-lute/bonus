package common

type Product struct {
	Id         int
	Title      string
	WebPath    string
	SmartLabel string
	Price      struct {
		Now struct {
			Amount float64
		}
		Was struct {
			Amount float64
		}
		UnitInfo struct {
			Price struct {
				Amount float64
			}
			Description string
		}
		Discount struct {
			SegmentId   int
			Description string
		}
	}
	Images []struct {
		Width    int
		Height   int
		Url      string
		TypeName string `graphql:"__typename"`
	}
}
